package grpc2fuse_test

import (
	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/stretchr/testify/mock"
)

type MockRawFileSystemClient struct {
	mock.Mock
	pb.RawFileSystemClient
}
