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

func (m *MockRawFileSystemClient) CopyFileRange(ctx context.Context, in *pb.CopyFileRangeRequest, opts ...grpc.CallOption) (*pb.CopyFileRangeResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.CopyFileRangeResponse), args.Error(1)
}

func TestFileSystem_CopyFileRange(t *testing.T) {
	tests := []struct {
		name        string
		input       *fuse.CopyFileRangeIn
		mockResp    *pb.CopyFileRangeResponse
		mockErr     error
		wantWritten uint32
		wantStatus  fuse.Status
	}{
		{
			name: "successful copy",
			input: &fuse.CopyFileRangeIn{
				InHeader:  fuse.InHeader{NodeId: 1},
				FhIn:      2,
				OffIn:     100,
				NodeIdOut: 3,
				FhOut:     4,
				OffOut:    200,
				Len:       1000,
				Flags:     0,
			},
			mockResp: &pb.CopyFileRangeResponse{
				Status:  &pb.Status{Code: 0},
				Written: 1000,
			},
			mockErr:     nil,
			wantWritten: 1000,
			wantStatus:  fuse.OK,
		},
		{
			name: "error response",
			input: &fuse.CopyFileRangeIn{
				InHeader: fuse.InHeader{NodeId: 1},
			},
			mockResp: &pb.CopyFileRangeResponse{
				Status:  &pb.Status{Code: int32(fuse.EACCES)},
				Written: 0,
			},
			mockErr:     nil,
			wantWritten: 0,
			wantStatus:  fuse.EACCES,
		},
		{
			name: "grpc error",
			input: &fuse.CopyFileRangeIn{
				InHeader: fuse.InHeader{NodeId: 1},
			},
			mockResp:    nil,
			mockErr:     context.DeadlineExceeded,
			wantWritten: 0,
			wantStatus:  fuse.EIO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockRawFileSystemClient)
			fs := &fileSystem{
				client: mockClient,
			}

			expectedRequest := &pb.CopyFileRangeRequest{
				Header:    toPbHeader(&tt.input.InHeader),
				FhIn:      tt.input.FhIn,
				OffIn:     tt.input.OffIn,
				NodeIdOut: tt.input.NodeIdOut,
				FhOut:     tt.input.FhOut,
				OffOut:    tt.input.OffOut,
				Len:       tt.input.Len,
				Flags:     tt.input.Flags,
			}

			mockClient.On("CopyFileRange", mock.Anything, expectedRequest, mock.Anything).Return(tt.mockResp, tt.mockErr)

			written, status := fs.CopyFileRange(make(chan struct{}), tt.input)

			assert.Equal(t, tt.wantWritten, written)
			assert.Equal(t, tt.wantStatus, status)
			mockClient.AssertExpectations(t)
		})
	}
}
