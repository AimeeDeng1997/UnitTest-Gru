package grpc2fuse

import (
	"context"
	"testing"

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

func (m *MockRawFileSystemClient) StatFs(ctx context.Context, in *pb.StatfsRequest, opts ...grpc.CallOption) (*pb.StatfsResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.StatfsResponse), args.Error(1)
}

func TestFileSystem_StatFs(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*MockRawFileSystemClient)
		input    *fuse.InHeader
		wantCode fuse.Status
		wantOut  *fuse.StatfsOut
	}{
		{
			name: "successful statfs",
			setup: func(m *MockRawFileSystemClient) {
				m.On("StatFs", mock.Anything, mock.MatchedBy(func(req *pb.StatfsRequest) bool {
					return req.Input.Length == 1 &&
						req.Input.Opcode == 2 &&
						req.Input.NodeId == 3
				}), mock.Anything).Return(&pb.StatfsResponse{
					Status: &pb.Status{Code: 0},
					Blocks: 1000,
					Bfree:  500,
					Bavail: 400,
					Files:  200,
					Ffree:  100,
					Bsize:  4096,
					NameLen: 255,
					Frsize: 4096,
					Padding: 0,
				}, nil)
			},
			input: &fuse.InHeader{
				Length:  1,
				Opcode: 2,
				NodeId: 3,
			},
			wantCode: fuse.OK,
			wantOut: &fuse.StatfsOut{
				Blocks:  1000,
				Bfree:   500,
				Bavail:  400,
				Files:   200,
				Ffree:   100,
				Bsize:   4096,
				NameLen: 255,
				Frsize:  4096,
				Padding: 0,
			},
		},
		{
			name: "failed statfs with error code",
			setup: func(m *MockRawFileSystemClient) {
				m.On("StatFs", mock.Anything, mock.Anything, mock.Anything).Return(&pb.StatfsResponse{
					Status: &pb.Status{Code: int32(fuse.ENOENT)},
				}, nil)
			},
			input: &fuse.InHeader{
				Length: 1,
			},
			wantCode: fuse.ENOENT,
			wantOut:  &fuse.StatfsOut{},
		},
		{
			name: "grpc error",
			setup: func(m *MockRawFileSystemClient) {
				m.On("StatFs", mock.Anything, mock.Anything, mock.Anything).Return(nil, grpc.ErrServerStopped)
			},
			input: &fuse.InHeader{
				Length: 1,
			},
			wantCode: fuse.EIO,
			wantOut:  &fuse.StatfsOut{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockRawFileSystemClient{}
			if tt.setup != nil {
				tt.setup(mockClient)
			}

			fs := &fileSystem{
				client: mockClient,
			}

			out := &fuse.StatfsOut{}
			cancel := make(chan struct{})

			got := fs.StatFs(cancel, tt.input, out)

			assert.Equal(t, tt.wantCode, got)
			if tt.wantCode == fuse.OK {
				assert.Equal(t, tt.wantOut, out)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
