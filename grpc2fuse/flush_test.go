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

// MockRawFileSystemClient is a mock for RawFileSystemClient
type MockRawFileSystemClient struct {
	mock.Mock
	pb.RawFileSystemClient
}

func (m *MockRawFileSystemClient) Flush(ctx context.Context, in *pb.FlushRequest, opts ...grpc.CallOption) (*pb.FlushResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.FlushResponse), args.Error(1)
}

func TestFileSystem_Flush(t *testing.T) {
	tests := []struct {
		name    string
		input   *fuse.FlushIn
		mockRes *pb.FlushResponse
		mockErr error
		want    fuse.Status
	}{
		{
			name: "successful flush",
			input: &fuse.FlushIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Fh:        123,
				LockOwner: 456,
			},
			mockRes: &pb.FlushResponse{
				Status: &pb.Status{
					Code: 0,
				},
			},
			mockErr: nil,
			want:    fuse.OK,
		},
		{
			name: "flush with error",
			input: &fuse.FlushIn{
				InHeader: fuse.InHeader{
					NodeId: 2,
				},
				Fh:        789,
				LockOwner: 101,
			},
			mockRes: &pb.FlushResponse{
				Status: &pb.Status{
					Code: int32(fuse.EACCES),
				},
			},
			mockErr: nil,
			want:    fuse.EACCES,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockRawFileSystemClient)
			fs := &fileSystem{
				client: mockClient,
			}

			expectedReq := &pb.FlushRequest{
				Header:    toPbHeader(&tt.input.InHeader),
				Fh:        tt.input.Fh,
				Unused:    tt.input.Unused,
				Padding:   tt.input.Padding,
				LockOwner: tt.input.LockOwner,
			}

			mockClient.On("Flush", mock.Anything, expectedReq, mock.Anything).Return(tt.mockRes, tt.mockErr)

			cancel := make(chan struct{})
			got := fs.Flush(cancel, tt.input)

			assert.Equal(t, tt.want, got)
			mockClient.AssertExpectations(t)
		})
	}
}
