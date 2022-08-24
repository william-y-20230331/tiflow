// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package enginepb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ResourceManagerClient is the client API for ResourceManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ResourceManagerClient interface {
	CreateResource(ctx context.Context, in *CreateResourceRequest, opts ...grpc.CallOption) (*CreateResourceResponse, error)
	QueryResource(ctx context.Context, in *QueryResourceRequest, opts ...grpc.CallOption) (*QueryResourceResponse, error)
	// RemoveResource cleans up the metadata only of the resource.
	// The invoker should handle the actual cleaning up on its own.
	RemoveResource(ctx context.Context, in *RemoveResourceRequest, opts ...grpc.CallOption) (*RemoveResourceResponse, error)
}

type resourceManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewResourceManagerClient(cc grpc.ClientConnInterface) ResourceManagerClient {
	return &resourceManagerClient{cc}
}

func (c *resourceManagerClient) CreateResource(ctx context.Context, in *CreateResourceRequest, opts ...grpc.CallOption) (*CreateResourceResponse, error) {
	out := new(CreateResourceResponse)
	err := c.cc.Invoke(ctx, "/enginepb.ResourceManager/CreateResource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceManagerClient) QueryResource(ctx context.Context, in *QueryResourceRequest, opts ...grpc.CallOption) (*QueryResourceResponse, error) {
	out := new(QueryResourceResponse)
	err := c.cc.Invoke(ctx, "/enginepb.ResourceManager/QueryResource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *resourceManagerClient) RemoveResource(ctx context.Context, in *RemoveResourceRequest, opts ...grpc.CallOption) (*RemoveResourceResponse, error) {
	out := new(RemoveResourceResponse)
	err := c.cc.Invoke(ctx, "/enginepb.ResourceManager/RemoveResource", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ResourceManagerServer is the server API for ResourceManager service.
// All implementations should embed UnimplementedResourceManagerServer
// for forward compatibility
type ResourceManagerServer interface {
	CreateResource(context.Context, *CreateResourceRequest) (*CreateResourceResponse, error)
	QueryResource(context.Context, *QueryResourceRequest) (*QueryResourceResponse, error)
	// RemoveResource cleans up the metadata only of the resource.
	// The invoker should handle the actual cleaning up on its own.
	RemoveResource(context.Context, *RemoveResourceRequest) (*RemoveResourceResponse, error)
}

// UnimplementedResourceManagerServer should be embedded to have forward compatible implementations.
type UnimplementedResourceManagerServer struct {
}

func (UnimplementedResourceManagerServer) CreateResource(context.Context, *CreateResourceRequest) (*CreateResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateResource not implemented")
}
func (UnimplementedResourceManagerServer) QueryResource(context.Context, *QueryResourceRequest) (*QueryResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryResource not implemented")
}
func (UnimplementedResourceManagerServer) RemoveResource(context.Context, *RemoveResourceRequest) (*RemoveResourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveResource not implemented")
}

// UnsafeResourceManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ResourceManagerServer will
// result in compilation errors.
type UnsafeResourceManagerServer interface {
	mustEmbedUnimplementedResourceManagerServer()
}

func RegisterResourceManagerServer(s grpc.ServiceRegistrar, srv ResourceManagerServer) {
	s.RegisterService(&ResourceManager_ServiceDesc, srv)
}

func _ResourceManager_CreateResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceManagerServer).CreateResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/enginepb.ResourceManager/CreateResource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceManagerServer).CreateResource(ctx, req.(*CreateResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ResourceManager_QueryResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceManagerServer).QueryResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/enginepb.ResourceManager/QueryResource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceManagerServer).QueryResource(ctx, req.(*QueryResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ResourceManager_RemoveResource_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveResourceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ResourceManagerServer).RemoveResource(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/enginepb.ResourceManager/RemoveResource",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ResourceManagerServer).RemoveResource(ctx, req.(*RemoveResourceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ResourceManager_ServiceDesc is the grpc.ServiceDesc for ResourceManager service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ResourceManager_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "enginepb.ResourceManager",
	HandlerType: (*ResourceManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateResource",
			Handler:    _ResourceManager_CreateResource_Handler,
		},
		{
			MethodName: "QueryResource",
			Handler:    _ResourceManager_QueryResource_Handler,
		},
		{
			MethodName: "RemoveResource",
			Handler:    _ResourceManager_RemoveResource_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "engine/proto/resources.proto",
}