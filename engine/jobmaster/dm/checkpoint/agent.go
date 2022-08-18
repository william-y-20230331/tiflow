// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package checkpoint

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pingcap/tidb-tools/pkg/dbutil"
	dmconfig "github.com/pingcap/tiflow/dm/config"
	"github.com/pingcap/tiflow/dm/pkg/conn"
	"github.com/pingcap/tiflow/dm/pkg/cputil"
	"github.com/pingcap/tiflow/engine/framework"
	"github.com/pingcap/tiflow/engine/jobmaster/dm/config"
	"github.com/pingcap/tiflow/engine/jobmaster/dm/metadata"
	"go.uber.org/zap"
)

const (
	loadCheckpointTable = `CREATE TABLE IF NOT EXISTS %s (
		task_name varchar(255) NOT NULL,
		source_name varchar(255) NOT NULL,
		status varchar(10) NOT NULL DEFAULT 'init' COMMENT 'init,running,finished',
		PRIMARY KEY (task_name, source_name)
	);`
	syncCheckpointTable = `CREATE TABLE IF NOT EXISTS %s (
		id VARCHAR(32) NOT NULL,
		cp_schema VARCHAR(128) NOT NULL,
		cp_table VARCHAR(128) NOT NULL,
		binlog_name VARCHAR(128),
		binlog_pos INT UNSIGNED,
		binlog_gtid TEXT,
		exit_safe_binlog_name VARCHAR(128) DEFAULT '',
		exit_safe_binlog_pos INT UNSIGNED DEFAULT 0,
		exit_safe_binlog_gtid TEXT,
		table_info JSON NOT NULL,
		is_global BOOLEAN,
		create_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		update_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		UNIQUE KEY uk_id_schema_table (id, cp_schema, cp_table)
	)`
)

// NewCheckpointAgent is a method to create a new checkpoint agent
var NewCheckpointAgent = NewAgentImpl

// Agent defeins a checkpoint agent interface
type Agent interface {
	Create(ctx context.Context, cfg *config.JobCfg) error
	Remove(ctx context.Context, cfg *config.JobCfg) error
	IsFresh(ctx context.Context, workerType framework.WorkerType, task *metadata.Task) (bool, error)
}

// AgentImpl implements Agent
type AgentImpl struct {
	logger *zap.Logger
}

// NewAgentImpl creates a new AgentImpl instance
func NewAgentImpl(pLogger *zap.Logger) Agent {
	c := &AgentImpl{
		logger: pLogger.With(zap.String("component", "checkpoint_agent")),
	}
	return c
}

// Create implements Agent.Create
// Create will be called when the job is created/restarted/updated.
// We create checkpoint table in master rather than in workers,
// to avoid the annoying log of "table already exists",
// because one job only need to create one checkpoint table per unit.
// move these codes to tiflow later.
func (c *AgentImpl) Create(ctx context.Context, cfg *config.JobCfg) error {
	c.logger.Info("create checkpoint")
	db, err := conn.DefaultDBProvider.Apply(cfg.TargetDB)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := createMetaDatabase(ctx, cfg, db); err != nil {
		return err
	}

	if cfg.TaskMode != dmconfig.ModeIncrement {
		if err := createLoadCheckpointTable(ctx, cfg, db); err != nil {
			return err
		}
	}
	if cfg.TaskMode != dmconfig.ModeFull {
		if err := createSyncCheckpointTable(ctx, cfg, db); err != nil {
			return err
		}
	}
	return nil
}

// Remove implements Agent.Remove
func (c *AgentImpl) Remove(ctx context.Context, cfg *config.JobCfg) error {
	c.logger.Info("remove checkpoint")
	db, err := conn.DefaultDBProvider.Apply(cfg.TargetDB)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := dropLoadCheckpointTable(ctx, cfg, db); err != nil {
		return err
	}
	return dropSyncCheckpointTable(ctx, cfg, db)
}

// IsFresh implements Agent.IsFresh
func (c *AgentImpl) IsFresh(ctx context.Context, workerType framework.WorkerType, task *metadata.Task) (bool, error) {
	if workerType == framework.WorkerDMDump {
		return true, nil
	}

	db, err := conn.DefaultDBProvider.Apply(task.Cfg.TargetDB)
	if err != nil {
		return false, err
	}
	defer db.Close()

	if workerType == framework.WorkerDMLoad {
		return isLoadFresh(ctx, task.Cfg, db)
	}
	return isSyncFresh(ctx, task.Cfg, db)
}

func createMetaDatabase(ctx context.Context, cfg *config.JobCfg, db *conn.BaseDB) error {
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbutil.ColumnName(cfg.MetaSchema))
	_, err := db.DB.ExecContext(ctx, query)
	return err
}

func createLoadCheckpointTable(ctx context.Context, cfg *config.JobCfg, db *conn.BaseDB) error {
	_, err := db.DB.ExecContext(ctx, fmt.Sprintf(loadCheckpointTable, loadTableName(cfg)))
	return err
}

func createSyncCheckpointTable(ctx context.Context, cfg *config.JobCfg, db *conn.BaseDB) error {
	_, err := db.DB.ExecContext(ctx, fmt.Sprintf(syncCheckpointTable, syncTableName(cfg)))
	return err
}

func dropLoadCheckpointTable(ctx context.Context, cfg *config.JobCfg, db *conn.BaseDB) error {
	dropTable := "DROP TABLE IF EXISTS %s"
	_, err := db.DB.ExecContext(ctx, fmt.Sprintf(dropTable, loadTableName(cfg)))
	return err
}

func dropSyncCheckpointTable(ctx context.Context, cfg *config.JobCfg, db *conn.BaseDB) error {
	dropTable := "DROP TABLE IF EXISTS %s"
	if _, err := db.DB.ExecContext(ctx, fmt.Sprintf(dropTable, syncTableName(cfg))); err != nil {
		return err
	}
	// The following two would be better removed in the worker when destroy.
	if _, err := db.DB.ExecContext(ctx, fmt.Sprintf(dropTable, shardMetaName(cfg))); err != nil {
		return err
	}
	if _, err := db.DB.ExecContext(ctx, fmt.Sprintf(dropTable, onlineDDLName(cfg))); err != nil {
		return err
	}
	return nil
}

func loadTableName(cfg *config.JobCfg) string {
	return dbutil.TableName(cfg.MetaSchema, cputil.LightningCheckpoint(cfg.Name))
}

func syncTableName(cfg *config.JobCfg) string {
	return dbutil.TableName(cfg.MetaSchema, cputil.SyncerCheckpoint(cfg.Name))
}

func shardMetaName(cfg *config.JobCfg) string {
	return dbutil.TableName(cfg.MetaSchema, cputil.SyncerShardMeta(cfg.Name))
}

func onlineDDLName(cfg *config.JobCfg) string {
	return dbutil.TableName(cfg.MetaSchema, cputil.SyncerOnlineDDL(cfg.Name))
}

func isLoadFresh(ctx context.Context, taskCfg *config.TaskCfg, db *conn.BaseDB) (bool, error) {
	// nolint:gosec
	query := fmt.Sprintf("SELECT status FROM %s WHERE `task_name` = ? AND `source_name` = ?", loadTableName((*config.JobCfg)(taskCfg)))
	var status string
	err := db.DB.QueryRowContext(ctx, query, taskCfg.Name, taskCfg.Upstreams[0].SourceID).Scan(&status)
	switch {
	case err == nil:
		return status == "init", nil
	case err == sql.ErrNoRows:
		return true, nil
	default:
		return false, err
	}
}

func isSyncFresh(ctx context.Context, taskCfg *config.TaskCfg, db *conn.BaseDB) (bool, error) {
	// nolint:gosec
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE `id` = ? AND `is_global` = true", syncTableName((*config.JobCfg)(taskCfg)))
	var status string
	err := db.DB.QueryRowContext(ctx, query, taskCfg.Upstreams[0].SourceID).Scan(&status)
	switch {
	case err == nil:
		return false, nil
	case err == sql.ErrNoRows:
		return true, nil
	default:
		return false, err
	}
}