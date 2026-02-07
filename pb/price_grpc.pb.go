
package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

const (
	PriceService_GetCurrentPrice_FullMethodName       = "/price.PriceService/GetCurrentPrice"
	PriceService_SubscribePriceUpdates_FullMethodName = "/price.PriceService/SubscribePriceUpdates"
)

// PriceServiceClient is the client API for PriceService service.
type PriceServiceClient interface {
	GetCurrentPrice(ctx context.Context, in *GetCurrentPriceRequest, opts ...grpc.CallOption) (*PriceUpdate, error)
	SubscribePriceUpdates(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (PriceService_SubscribePriceUpdatesClient, error)
}

type priceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPriceServiceClient(cc grpc.ClientConnInterface) PriceServiceClient {
	return &priceServiceClient{cc}
}

func (c *priceServiceClient) GetCurrentPrice(ctx context.Context, in *GetCurrentPriceRequest, opts ...grpc.CallOption) (*PriceUpdate, error) {
	out := new(PriceUpdate)
	err := c.cc.Invoke(ctx, PriceService_GetCurrentPrice_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *priceServiceClient) SubscribePriceUpdates(ctx context.Context, in *SubscribeRequest, opts ...grpc.CallOption) (PriceService_SubscribePriceUpdatesClient, error) {
	stream, err := c.cc.NewStream(ctx, &PriceService_ServiceDesc.Streams[0], PriceService_SubscribePriceUpdates_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &priceServiceSubscribePriceUpdatesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PriceService_SubscribePriceUpdatesClient interface {
	Recv() (*PriceUpdate, error)
	grpc.ClientStream
}

type priceServiceSubscribePriceUpdatesClient struct {
	grpc.ClientStream
}

func (x *priceServiceSubscribePriceUpdatesClient) Recv() (*PriceUpdate, error) {
	m := new(PriceUpdate)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PriceServiceServer is the server API for PriceService service.
type PriceServiceServer interface {
	GetCurrentPrice(context.Context, *GetCurrentPriceRequest) (*PriceUpdate, error)
	SubscribePriceUpdates(*SubscribeRequest, PriceService_SubscribePriceUpdatesServer) error
}

// UnimplementedPriceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPriceServiceServer struct{}

func (UnimplementedPriceServiceServer) GetCurrentPrice(context.Context, *GetCurrentPriceRequest) (*PriceUpdate, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCurrentPrice not implemented")
}
func (UnimplementedPriceServiceServer) SubscribePriceUpdates(*SubscribeRequest, PriceService_SubscribePriceUpdatesServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribePriceUpdates not implemented")
}

// UnsafePriceServiceServer may be embedded to opt out of forward compatibility for this service.
type UnsafePriceServiceServer interface {
	mustEmbedUnimplementedPriceServiceServer()
}

func RegisterPriceServiceServer(s grpc.ServiceRegistrar, srv PriceServiceServer) {
	s.RegisterService(&PriceService_ServiceDesc, srv)
}

func _PriceService_GetCurrentPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCurrentPriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PriceServiceServer).GetCurrentPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PriceService_GetCurrentPrice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PriceServiceServer).GetCurrentPrice(ctx, req.(*GetCurrentPriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PriceService_SubscribePriceUpdates_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PriceServiceServer).SubscribePriceUpdates(m, &priceServiceSubscribePriceUpdatesServer{stream})
}

type PriceService_SubscribePriceUpdatesServer interface {
	Send(*PriceUpdate) error
	grpc.ServerStream
}

type priceServiceSubscribePriceUpdatesServer struct {
	grpc.ServerStream
}

func (x *priceServiceSubscribePriceUpdatesServer) Send(m *PriceUpdate) error {
	return x.ServerStream.SendMsg(m)
}

// PriceService_ServiceDesc is the grpc.ServiceDesc for PriceService service.
var PriceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "price.PriceService",
	HandlerType: (*PriceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCurrentPrice",
			Handler:    _PriceService_GetCurrentPrice_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribePriceUpdates",
			Handler:       _PriceService_SubscribePriceUpdates_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/price.proto",
}
