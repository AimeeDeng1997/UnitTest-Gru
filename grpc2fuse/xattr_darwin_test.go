package grpc2fuse_test

import (
	"context"
	"testing"

	"github.com/chiyutianyi/grpcfuse/grpc2fuse"
	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type mockRawFileSystemClient struct {
	pb.RawFileSystemClient
	setXAttrFunc func(ctx context.Context, in *pb.SetXAttrRequest, opts ...grpc.CallOption) (*pb.SetXAttrResponse, error)
}

func (m *mockRawFileSystemClient) SetXAttr(ctx context.Context, in *pb.SetXAttrRequest, opts ...grpc.CallOption) (*pb.SetXAttrResponse, error) {
	if m.setXAttrFunc != nil {
		return m.setXAttrFunc(ctx, in, opts...)
	}
	return nil, nil
}

func TestSetXAttr(t *testing.T) {
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
				Flags:    1,
				Position: 0,
				Padding:  0,
			},
			attr: "user.test",
			data: []byte("test value"),
			mockResp: &pb.SetXAttrResponse{
				Status: &pb.Status{
					Code: 0,
				},
			},
			want: fuse.OK,
		},
		{
			name: "error status from server",
			input: &fuse.SetXAttrIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Size:  10,
				Flags: 1,
			},
			attr: "user.test",
			data: []byte("test value"),
			mockResp: &pb.SetXAttrResponse{
				Status: &pb.Status{
					Code: int32(fuse.EACCES),
				},
			},
			want: fuse.EACCES,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockRawFileSystemClient{
				setXAttrFunc: func(ctx context.Context, in *pb.SetXAttrRequest, opts ...grpc.CallOption) (*pb.SetXAttrResponse, error) {
					assert.Equal(t, tt.input.NodeId, in.Header.NodeId)
					assert.Equal(t, tt.attr, in.Attr)
					assert.Equal(t, tt.data, in.Data)
					assert.Equal(t, tt.input.Size, in.Size)
					assert.Equal(t, tt.input.Flags, in.Flags)
					assert.Equal(t, tt.input.Position, in.Position)
					assert.Equal(t, tt.input.Padding, in.Padding)
					return tt.mockResp, tt.mockErr
				},
			}

			fs := grpc2fuse.NewFileSystem(mockClient)
			cancel := make(chan struct{})

			got := fs.SetXAttr(cancel, tt.input, tt.attr, tt.data)
			assert.Equal(t, tt.want, got)
		})
	}
}
