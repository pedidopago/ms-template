// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package xyzv1

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

// XYZServiceClient is the client API for XYZService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type XYZServiceClient interface {
	Example(ctx context.Context, in *ExampleRequest, opts ...grpc.CallOption) (*ExampleResponse, error)
}

type xYZServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewXYZServiceClient(cc grpc.ClientConnInterface) XYZServiceClient {
	return &xYZServiceClient{cc}
}

func (c *xYZServiceClient) Example(ctx context.Context, in *ExampleRequest, opts ...grpc.CallOption) (*ExampleResponse, error) {
	out := new(ExampleResponse)
	err := c.cc.Invoke(ctx, "/pedidopago.xyz.v1.XYZService/Example", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// XYZServiceServer is the server API for XYZService service.
// All implementations should embed UnimplementedXYZServiceServer
// for forward compatibility
type XYZServiceServer interface {
	Example(context.Context, *ExampleRequest) (*ExampleResponse, error)
}

// UnimplementedXYZServiceServer should be embedded to have forward compatible implementations.
type UnimplementedXYZServiceServer struct {
}

func (UnimplementedXYZServiceServer) Example(context.Context, *ExampleRequest) (*ExampleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Example not implemented")
}

// UnsafeXYZServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to XYZServiceServer will
// result in compilation errors.
type UnsafeXYZServiceServer interface {
	mustEmbedUnimplementedXYZServiceServer()
}

func RegisterXYZServiceServer(s grpc.ServiceRegistrar, srv XYZServiceServer) {
	s.RegisterService(&XYZService_ServiceDesc, srv)
}

func _XYZService_Example_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExampleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(XYZServiceServer).Example(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pedidopago.xyz.v1.XYZService/Example",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(XYZServiceServer).Example(ctx, req.(*ExampleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// XYZService_ServiceDesc is the grpc.ServiceDesc for XYZService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var XYZService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pedidopago.xyz.v1.XYZService",
	HandlerType: (*XYZServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Example",
			Handler:    _XYZService_Example_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pedidopago/xyz/v1/xyz_service.proto",
}
