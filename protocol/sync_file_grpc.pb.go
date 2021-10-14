// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protocol

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

// SyncFileClient is the client API for SyncFile service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SyncFileClient interface {
	SyncFile(ctx context.Context, in *SyncFileRequest, opts ...grpc.CallOption) (*Result, error)
	SyncRepoFile(ctx context.Context, in *SyncRepoFileRequest, opts ...grpc.CallOption) (*Result, error)
}

type syncFileClient struct {
	cc grpc.ClientConnInterface
}

func NewSyncFileClient(cc grpc.ClientConnInterface) SyncFileClient {
	return &syncFileClient{cc}
}

func (c *syncFileClient) SyncFile(ctx context.Context, in *SyncFileRequest, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/syncfile.SyncFile/SyncFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *syncFileClient) SyncRepoFile(ctx context.Context, in *SyncRepoFileRequest, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/syncfile.SyncFile/SyncRepoFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SyncFileServer is the server API for SyncFile service.
// All implementations must embed UnimplementedSyncFileServer
// for forward compatibility
type SyncFileServer interface {
	SyncFile(context.Context, *SyncFileRequest) (*Result, error)
	SyncRepoFile(context.Context, *SyncRepoFileRequest) (*Result, error)
	mustEmbedUnimplementedSyncFileServer()
}

// UnimplementedSyncFileServer must be embedded to have forward compatible implementations.
type UnimplementedSyncFileServer struct {
}

func (UnimplementedSyncFileServer) SyncFile(context.Context, *SyncFileRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncFile not implemented")
}
func (UnimplementedSyncFileServer) SyncRepoFile(context.Context, *SyncRepoFileRequest) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncRepoFile not implemented")
}
func (UnimplementedSyncFileServer) mustEmbedUnimplementedSyncFileServer() {}

// UnsafeSyncFileServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SyncFileServer will
// result in compilation errors.
type UnsafeSyncFileServer interface {
	mustEmbedUnimplementedSyncFileServer()
}

func RegisterSyncFileServer(s grpc.ServiceRegistrar, srv SyncFileServer) {
	s.RegisterService(&SyncFile_ServiceDesc, srv)
}

func _SyncFile_SyncFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SyncFileServer).SyncFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syncfile.SyncFile/SyncFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SyncFileServer).SyncFile(ctx, req.(*SyncFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SyncFile_SyncRepoFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncRepoFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SyncFileServer).SyncRepoFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/syncfile.SyncFile/SyncRepoFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SyncFileServer).SyncRepoFile(ctx, req.(*SyncRepoFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SyncFile_ServiceDesc is the grpc.ServiceDesc for SyncFile service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SyncFile_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "syncfile.SyncFile",
	HandlerType: (*SyncFileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SyncFile",
			Handler:    _SyncFile_SyncFile_Handler,
		},
		{
			MethodName: "SyncRepoFile",
			Handler:    _SyncFile_SyncRepoFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sync_file.proto",
}