package fuse2grpc_test

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/chiyutianyi/grpcfuse/fuse2grpc"
	"github.com/chiyutianyi/grpcfuse/pb"
)

type MockRawFileSystem struct {
	fuse.RawFileSystem
}

func (m *MockRawFileSystem) SetXAttr(cancel <-chan struct{}, in *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
	return fuse.OK
}

func (m *MockRawFileSystem) String() string {
	return "MockRawFileSystem"
}

func TestSetXAttr_Success(t *testing.T) {
	mockFS := &MockRawFileSystem{}
	server := fuse2grpc.NewServer(mockFS)

	req := &pb.SetXAttrRequest{
		Header: &pb.InHeader{
			NodeId: 1,
			Caller: &pb.Caller{
				Owner: &pb.Owner{
					Uid: 1000,
					Gid: 1000,
				},
			},
		},
		Attr: "user.test",
		Size: 4,
		Data: []byte("test"),
	}

	resp, err := server.SetXAttr(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, &pb.SetXAttrResponse{
		Status: &pb.Status{Code: 0},
	}, resp)
}
