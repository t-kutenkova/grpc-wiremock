// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.18.1
// source: api/grpc/wearable.proto

package wearable

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

// WearableServiceClient is the client API for WearableService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WearableServiceClient interface {
	BeatsPerMinute(ctx context.Context, in *BeatsPerMinuteRequest, opts ...grpc.CallOption) (WearableService_BeatsPerMinuteClient, error)
	ConsumeBeatsPerMinute(ctx context.Context, opts ...grpc.CallOption) (WearableService_ConsumeBeatsPerMinuteClient, error)
	CalculateBeatsPerMinute(ctx context.Context, opts ...grpc.CallOption) (WearableService_CalculateBeatsPerMinuteClient, error)
}

type wearableServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWearableServiceClient(cc grpc.ClientConnInterface) WearableServiceClient {
	return &wearableServiceClient{cc}
}

func (c *wearableServiceClient) BeatsPerMinute(ctx context.Context, in *BeatsPerMinuteRequest, opts ...grpc.CallOption) (WearableService_BeatsPerMinuteClient, error) {
	stream, err := c.cc.NewStream(ctx, &WearableService_ServiceDesc.Streams[0], "/wearable.WearableService/BeatsPerMinute", opts...)
	if err != nil {
		return nil, err
	}
	x := &wearableServiceBeatsPerMinuteClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type WearableService_BeatsPerMinuteClient interface {
	Recv() (*BeatsPerMinuteResponse, error)
	grpc.ClientStream
}

type wearableServiceBeatsPerMinuteClient struct {
	grpc.ClientStream
}

func (x *wearableServiceBeatsPerMinuteClient) Recv() (*BeatsPerMinuteResponse, error) {
	m := new(BeatsPerMinuteResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *wearableServiceClient) ConsumeBeatsPerMinute(ctx context.Context, opts ...grpc.CallOption) (WearableService_ConsumeBeatsPerMinuteClient, error) {
	stream, err := c.cc.NewStream(ctx, &WearableService_ServiceDesc.Streams[1], "/wearable.WearableService/ConsumeBeatsPerMinute", opts...)
	if err != nil {
		return nil, err
	}
	x := &wearableServiceConsumeBeatsPerMinuteClient{stream}
	return x, nil
}

type WearableService_ConsumeBeatsPerMinuteClient interface {
	Send(*ConsumeBeatsPerMinuteRequest) error
	CloseAndRecv() (*ConsumeBeatsPerMinuteResponse, error)
	grpc.ClientStream
}

type wearableServiceConsumeBeatsPerMinuteClient struct {
	grpc.ClientStream
}

func (x *wearableServiceConsumeBeatsPerMinuteClient) Send(m *ConsumeBeatsPerMinuteRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *wearableServiceConsumeBeatsPerMinuteClient) CloseAndRecv() (*ConsumeBeatsPerMinuteResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(ConsumeBeatsPerMinuteResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *wearableServiceClient) CalculateBeatsPerMinute(ctx context.Context, opts ...grpc.CallOption) (WearableService_CalculateBeatsPerMinuteClient, error) {
	stream, err := c.cc.NewStream(ctx, &WearableService_ServiceDesc.Streams[2], "/wearable.WearableService/CalculateBeatsPerMinute", opts...)
	if err != nil {
		return nil, err
	}
	x := &wearableServiceCalculateBeatsPerMinuteClient{stream}
	return x, nil
}

type WearableService_CalculateBeatsPerMinuteClient interface {
	Send(*CalculateBeatsPerMinuteRequest) error
	Recv() (*CalculateBeatsPerMinuteResponse, error)
	grpc.ClientStream
}

type wearableServiceCalculateBeatsPerMinuteClient struct {
	grpc.ClientStream
}

func (x *wearableServiceCalculateBeatsPerMinuteClient) Send(m *CalculateBeatsPerMinuteRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *wearableServiceCalculateBeatsPerMinuteClient) Recv() (*CalculateBeatsPerMinuteResponse, error) {
	m := new(CalculateBeatsPerMinuteResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// WearableServiceServer is the server API for WearableService service.
// All implementations must embed UnimplementedWearableServiceServer
// for forward compatibility
type WearableServiceServer interface {
	BeatsPerMinute(*BeatsPerMinuteRequest, WearableService_BeatsPerMinuteServer) error
	ConsumeBeatsPerMinute(WearableService_ConsumeBeatsPerMinuteServer) error
	CalculateBeatsPerMinute(WearableService_CalculateBeatsPerMinuteServer) error
	mustEmbedUnimplementedWearableServiceServer()
}

// UnimplementedWearableServiceServer must be embedded to have forward compatible implementations.
type UnimplementedWearableServiceServer struct {
}

func (UnimplementedWearableServiceServer) BeatsPerMinute(*BeatsPerMinuteRequest, WearableService_BeatsPerMinuteServer) error {
	return status.Errorf(codes.Unimplemented, "method BeatsPerMinute not implemented")
}
func (UnimplementedWearableServiceServer) ConsumeBeatsPerMinute(WearableService_ConsumeBeatsPerMinuteServer) error {
	return status.Errorf(codes.Unimplemented, "method ConsumeBeatsPerMinute not implemented")
}
func (UnimplementedWearableServiceServer) CalculateBeatsPerMinute(WearableService_CalculateBeatsPerMinuteServer) error {
	return status.Errorf(codes.Unimplemented, "method CalculateBeatsPerMinute not implemented")
}
func (UnimplementedWearableServiceServer) mustEmbedUnimplementedWearableServiceServer() {}

// UnsafeWearableServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WearableServiceServer will
// result in compilation errors.
type UnsafeWearableServiceServer interface {
	mustEmbedUnimplementedWearableServiceServer()
}

func RegisterWearableServiceServer(s grpc.ServiceRegistrar, srv WearableServiceServer) {
	s.RegisterService(&WearableService_ServiceDesc, srv)
}

func _WearableService_BeatsPerMinute_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(BeatsPerMinuteRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(WearableServiceServer).BeatsPerMinute(m, &wearableServiceBeatsPerMinuteServer{stream})
}

type WearableService_BeatsPerMinuteServer interface {
	Send(*BeatsPerMinuteResponse) error
	grpc.ServerStream
}

type wearableServiceBeatsPerMinuteServer struct {
	grpc.ServerStream
}

func (x *wearableServiceBeatsPerMinuteServer) Send(m *BeatsPerMinuteResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _WearableService_ConsumeBeatsPerMinute_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(WearableServiceServer).ConsumeBeatsPerMinute(&wearableServiceConsumeBeatsPerMinuteServer{stream})
}

type WearableService_ConsumeBeatsPerMinuteServer interface {
	SendAndClose(*ConsumeBeatsPerMinuteResponse) error
	Recv() (*ConsumeBeatsPerMinuteRequest, error)
	grpc.ServerStream
}

type wearableServiceConsumeBeatsPerMinuteServer struct {
	grpc.ServerStream
}

func (x *wearableServiceConsumeBeatsPerMinuteServer) SendAndClose(m *ConsumeBeatsPerMinuteResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *wearableServiceConsumeBeatsPerMinuteServer) Recv() (*ConsumeBeatsPerMinuteRequest, error) {
	m := new(ConsumeBeatsPerMinuteRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _WearableService_CalculateBeatsPerMinute_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(WearableServiceServer).CalculateBeatsPerMinute(&wearableServiceCalculateBeatsPerMinuteServer{stream})
}

type WearableService_CalculateBeatsPerMinuteServer interface {
	Send(*CalculateBeatsPerMinuteResponse) error
	Recv() (*CalculateBeatsPerMinuteRequest, error)
	grpc.ServerStream
}

type wearableServiceCalculateBeatsPerMinuteServer struct {
	grpc.ServerStream
}

func (x *wearableServiceCalculateBeatsPerMinuteServer) Send(m *CalculateBeatsPerMinuteResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *wearableServiceCalculateBeatsPerMinuteServer) Recv() (*CalculateBeatsPerMinuteRequest, error) {
	m := new(CalculateBeatsPerMinuteRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// WearableService_ServiceDesc is the grpc.ServiceDesc for WearableService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WearableService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "wearable.WearableService",
	HandlerType: (*WearableServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BeatsPerMinute",
			Handler:       _WearableService_BeatsPerMinute_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ConsumeBeatsPerMinute",
			Handler:       _WearableService_ConsumeBeatsPerMinute_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "CalculateBeatsPerMinute",
			Handler:       _WearableService_CalculateBeatsPerMinute_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api/grpc/wearable.proto",
}
