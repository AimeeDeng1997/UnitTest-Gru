package grpc2fuse

import (
	"context"
	"errors"
	"testing"

	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Mock RawFileSystemClient
type mockRawFileSystemClient struct {
	mock.Mock
	pb.RawFileSystemClient
}

func (m *mockRawFileSystemClient) Release(ctx context.Context, in *pb.ReleaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func TestFileSystem_Release(t *testing.T) {
	tests := []struct {
		name    string
		in      *fuse.ReleaseIn
		wantErr bool
		mockFn  func(*mockRawFileSystemClient)
	}{
		{
			name: "successful release",
			in: &fuse.ReleaseIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Fh:           123,
				Flags:        0,
				ReleaseFlags: 0,
				LockOwner:    0,
			},
			wantErr: false,
			mockFn: func(m *mockRawFileSystemClient) {
				m.On("Release", mock.Anything, mock.MatchedBy(func(req *pb.ReleaseRequest) bool {
					return req.Header.NodeId == 1 &&
						req.Fh == 123 &&
						req.Flags == 0 &&
						req.ReleaseFlags == 0 &&
						req.LockOwner == 0
				}), mock.Anything).Return(&emptypb.Empty{}, nil)
			},
		},
		{
			name: "release with error",
			in: &fuse.ReleaseIn{
				InHeader: fuse.InHeader{
					NodeId: 2,
				},
				Fh: 456,
			},
			wantErr: true,
			mockFn: func(m *mockRawFileSystemClient) {
				m.On("Release", mock.Anything, mock.Anything, mock.Anything).Return(&emptypb.Empty{}, errors.New("release error"))
			},
		},
		{
			name: "release with different flags",
			in: &fuse.ReleaseIn{
				InHeader: fuse.InHeader{
					NodeId: 3,
				},
				Fh:           789,
				Flags:        1,
				ReleaseFlags: 1,
				LockOwner:    1000,
			},
			wantErr: false,
			mockFn: func(m *mockRawFileSystemClient) {
				m.On("Release", mock.Anything, mock.MatchedBy(func(req *pb.ReleaseRequest) bool {
					return req.Header.NodeId == 3 &&
						req.Fh == 789 &&
						req.Flags == 1 &&
						req.ReleaseFlags == 1 &&
						req.LockOwner == 1000
				}), mock.Anything).Return(&emptypb.Empty{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockRawFileSystemClient{}
			if tt.mockFn != nil {
				tt.mockFn(mockClient)
			}

			fs := &fileSystem{
				client: mockClient,
			}

			cancel := make(chan struct{})
			fs.Release(cancel, tt.in)

			mockClient.AssertExpectations(t)
		})
	}
}
