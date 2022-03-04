// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: sf/ethereum/trxstream/v1/trxstream.proto

package pbtrxstream

import (
	context "context"
	v1 "github.com/streamingfast/sf-ethereum/pb/sf/ethereum/codec/v1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TransactionStreamClient is the client API for TransactionStream service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TransactionStreamClient interface {
	Transactions(ctx context.Context, in *TransactionRequest, opts ...grpc.CallOption) (TransactionStream_TransactionsClient, error)
}

type transactionStreamClient struct {
	cc grpc.ClientConnInterface
}

func NewTransactionStreamClient(cc grpc.ClientConnInterface) TransactionStreamClient {
	return &transactionStreamClient{cc}
}

func (c *transactionStreamClient) Transactions(ctx context.Context, in *TransactionRequest, opts ...grpc.CallOption) (TransactionStream_TransactionsClient, error) {
	stream, err := c.cc.NewStream(ctx, &TransactionStream_ServiceDesc.Streams[0], "/sf.ethereum.trxstream.v1.TransactionStream/Transactions", opts...)
	if err != nil {
		return nil, err
	}
	x := &transactionStreamTransactionsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TransactionStream_TransactionsClient interface {
	Recv() (*v1.Transaction, error)
	grpc.ClientStream
}

type transactionStreamTransactionsClient struct {
	grpc.ClientStream
}

func (x *transactionStreamTransactionsClient) Recv() (*v1.Transaction, error) {
	m := new(v1.Transaction)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TransactionStreamServer is the server API for TransactionStream service.
// All implementations should embed UnimplementedTransactionStreamServer
// for forward compatibility
type TransactionStreamServer interface {
	Transactions(*TransactionRequest, TransactionStream_TransactionsServer) error
}

// UnimplementedTransactionStreamServer should be embedded to have forward compatible implementations.
type UnimplementedTransactionStreamServer struct {
}

func (UnimplementedTransactionStreamServer) Transactions(*TransactionRequest, TransactionStream_TransactionsServer) error {
	return status.Errorf(codes.Unimplemented, "method Transactions not implemented")
}

// UnsafeTransactionStreamServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TransactionStreamServer will
// result in compilation errors.
type UnsafeTransactionStreamServer interface {
	mustEmbedUnimplementedTransactionStreamServer()
}

func RegisterTransactionStreamServer(s grpc.ServiceRegistrar, srv TransactionStreamServer) {
	s.RegisterService(&TransactionStream_ServiceDesc, srv)
}

func _TransactionStream_Transactions_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TransactionRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TransactionStreamServer).Transactions(m, &transactionStreamTransactionsServer{stream})
}

type TransactionStream_TransactionsServer interface {
	Send(*v1.Transaction) error
	grpc.ServerStream
}

type transactionStreamTransactionsServer struct {
	grpc.ServerStream
}

func (x *transactionStreamTransactionsServer) Send(m *v1.Transaction) error {
	return x.ServerStream.SendMsg(m)
}

// TransactionStream_ServiceDesc is the grpc.ServiceDesc for TransactionStream service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TransactionStream_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sf.ethereum.trxstream.v1.TransactionStream",
	HandlerType: (*TransactionStreamServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Transactions",
			Handler:       _TransactionStream_Transactions_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "sf/ethereum/trxstream/v1/trxstream.proto",
}
