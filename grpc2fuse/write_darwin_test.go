package grpc2fuse

import (
	"context"
	"testing"

	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type mockWriteClient struct {
	pb.RawFileSystemClient
	writeFunc func(context.Context, *pb.WriteRequest, ...grpc.CallOption) (*pb.WriteResponse, error)
}

func (m *mockWriteClient) Write(ctx context.Context, req *pb.WriteRequest, opts ...grpc.CallOption) (*pb.WriteResponse, error) {
	if m.writeFunc != nil {
		return m.writeFunc(ctx, req, opts...)
	}
	return nil, nil
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name        string
		input       *fuse.WriteIn
		data        []byte
		mockResp    *pb.WriteResponse
		mockErr     error
		wantWritten uint32
		wantStatus  fuse.Status
	}{
		{
			name: "successful write",
			input: &fuse.WriteIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Fh:         2,
				Offset:     100,
				Size:       5,
				WriteFlags: 0,
			},
			data: []byte("hello"),
			mockResp: &pb.WriteResponse{
				Status:  &pb.Status{Code: 0},
				Written: 5,
			},
			wantWritten: 5,
			wantStatus:  fuse.OK,
		},
		{
			name: "write error",
			input: &fuse.WriteIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Fh:         2,
				Offset:     0,
				Size:       3,
				WriteFlags: 0,
			},
			data: []byte("abc"),
			mockResp: &pb.WriteResponse{
				Status:  &pb.Status{Code: int32(fuse.EIO)},
				Written: 0,
			},
			wantWritten: 0,
			wantStatus:  fuse.EIO,
		},
		{
			name: "zero byte write",
			input: &fuse.WriteIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Fh:         2,
				Offset:     0,
				Size:       0,
				WriteFlags: 0,
			},
			data: []byte{},
			mockResp: &pb.WriteResponse{
				Status:  &pb.Status{Code: 0},
				Written: 0,
			},
			wantWritten: 0,
			wantStatus:  fuse.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockWriteClient{
				writeFunc: func(ctx context.Context, req *pb.WriteRequest, opts ...grpc.CallOption) (*pb.WriteResponse, error) {
					assert.Equal(t, tt.input.NodeId, req.Header.NodeId)
					assert.Equal(t, tt.input.Fh, req.Fh)
					assert.Equal(t, tt.input.Offset, req.Offset)
					assert.Equal(t, tt.input.Size, req.Size)
					assert.Equal(t, tt.input.WriteFlags, req.WriteFlags)
					assert.Equal(t, tt.data, req.Data)
					return tt.mockResp, tt.mockErr
				},
			}

			fs := &fileSystem{
				client: mockClient,
			}

			written, status := fs.Write(make(<-chan struct{}), tt.input, tt.data)
			assert.Equal(t, tt.wantWritten, written)
			assert.Equal(t, tt.wantStatus, status)
		})
	}
}
