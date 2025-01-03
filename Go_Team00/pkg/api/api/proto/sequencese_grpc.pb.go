// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.1
// source: api/proto/sequencese.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	FrequencyServise_GenerateFrequency_FullMethodName = "/api.FrequencyServise/GenerateFrequency"
)

// FrequencyServiseClient is the client API for FrequencyServise service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FrequencyServiseClient interface {
	GenerateFrequency(ctx context.Context, in *Frequency, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Frequency], error)
}

type frequencyServiseClient struct {
	cc grpc.ClientConnInterface
}

func NewFrequencyServiseClient(cc grpc.ClientConnInterface) FrequencyServiseClient {
	return &frequencyServiseClient{cc}
}

func (c *frequencyServiseClient) GenerateFrequency(ctx context.Context, in *Frequency, opts ...grpc.CallOption) (grpc.ServerStreamingClient[Frequency], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &FrequencyServise_ServiceDesc.Streams[0], FrequencyServise_GenerateFrequency_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Frequency, Frequency]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FrequencyServise_GenerateFrequencyClient = grpc.ServerStreamingClient[Frequency]

// FrequencyServiseServer is the server API for FrequencyServise service.
// All implementations must embed UnimplementedFrequencyServiseServer
// for forward compatibility.
type FrequencyServiseServer interface {
	GenerateFrequency(*Frequency, grpc.ServerStreamingServer[Frequency]) error
	mustEmbedUnimplementedFrequencyServiseServer()
}

// UnimplementedFrequencyServiseServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedFrequencyServiseServer struct{}

func (UnimplementedFrequencyServiseServer) GenerateFrequency(*Frequency, grpc.ServerStreamingServer[Frequency]) error {
	return status.Errorf(codes.Unimplemented, "method GenerateFrequency not implemented")
}
func (UnimplementedFrequencyServiseServer) mustEmbedUnimplementedFrequencyServiseServer() {}
func (UnimplementedFrequencyServiseServer) testEmbeddedByValue()                          {}

// UnsafeFrequencyServiseServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FrequencyServiseServer will
// result in compilation errors.
type UnsafeFrequencyServiseServer interface {
	mustEmbedUnimplementedFrequencyServiseServer()
}

func RegisterFrequencyServiseServer(s grpc.ServiceRegistrar, srv FrequencyServiseServer) {
	// If the following call pancis, it indicates UnimplementedFrequencyServiseServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&FrequencyServise_ServiceDesc, srv)
}

func _FrequencyServise_GenerateFrequency_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Frequency)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FrequencyServiseServer).GenerateFrequency(m, &grpc.GenericServerStream[Frequency, Frequency]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type FrequencyServise_GenerateFrequencyServer = grpc.ServerStreamingServer[Frequency]

// FrequencyServise_ServiceDesc is the grpc.ServiceDesc for FrequencyServise service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FrequencyServise_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.FrequencyServise",
	HandlerType: (*FrequencyServiseServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GenerateFrequency",
			Handler:       _FrequencyServise_GenerateFrequency_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/proto/sequencese.proto",
}
