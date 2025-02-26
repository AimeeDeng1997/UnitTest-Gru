package grpc2fuse_test

import (
	"context"
	"testing"

	"github.com/chiyutianyi/grpcfuse/grpc2fuse"
	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockRawFileSystemClient struct {
	mock.Mock
	pb.RawFileSystemClient
}

func (m *MockRawFileSystemClient) SetXAttr(ctx context.Context, in *pb.SetXAttrRequest, opts ...grpc.CallOption) (*pb.SetXAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.SetXAttrResponse), args.Error(1)
}

func TestSetXAttr(t *testing.T) {
	mockClient := new(MockRawFileSystemClient)
	fs := grpc2fuse.NewFileSystem(mockClient)

	tests := []struct {
		name     string
		input    *fuse.SetXAttrIn
		attr     string
		data     []byte
		mockResp *pb.SetXAttrResponse
		mockErr  error
		want     fuse.Status
	}{
		{
			name: "successful set xattr",
			input: &fuse.SetXAttrIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Size:     10,
				Flags:    0,
				Position: 0,
				Padding:  0,
			},
			attr: "user.test",
			data: []byte("test value"),
			mockResp: &pb.SetXAttrResponse{
				Status: &pb.Status{Code: 0},
			},
			mockErr: nil,
			want:    fuse.OK,
		},
		{
			name: "error setting xattr",
			input: &fuse.SetXAttrIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Size:     10,
				Flags:    0,
				Position: 0,
				Padding:  0,
			},
			attr: "user.test",
			data: []byte("test value"),
			mockResp: &pb.SetXAttrResponse{
				Status: &pb.Status{Code: int32(fuse.EACCES)},
			},
			mockErr: nil,
			want:    fuse.EACCES,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("SetXAttr", mock.Anything, &pb.SetXAttrRequest{
				Header: &pb.InHeader{
					NodeId: tt.input.NodeId,
				},
				Attr:     tt.attr,
				Data:     tt.data,
				Size:     tt.input.Size,
				Flags:    tt.input.Flags,
				Position: tt.input.Position,
				Padding:  tt.input.Padding,
			}, mock.Anything).Return(tt.mockResp, tt.mockErr)

			cancel := make(chan struct{})
			got := fs.SetXAttr(cancel, tt.input, tt.attr, tt.data)
			assert.Equal(t, tt.want, got)
			mockClient.AssertExpectations(t)
		})
	}
}
