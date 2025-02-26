package fuse2grpc

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/pb"
)

type mockRawFileSystem struct {
	mock.Mock
	fuse.RawFileSystem
}

func (m *mockRawFileSystem) Mknod(cancel <-chan struct{}, in *fuse.MknodIn, name string, out *fuse.EntryOut) fuse.Status {
	args := m.Called(cancel, in, name, out)
	return args.Get(0).(fuse.Status)
}

func (m *mockRawFileSystem) String() string {
	return "mockRawFileSystem"
}

func TestMknod(t *testing.T) {
	tests := []struct {
		name        string
		req         *pb.MknodRequest
		setupMock   func(*mockRawFileSystem)
		wantErr     bool
		wantErrCode codes.Code
		wantStatus  int32
	}{
		{
			name: "success",
			req: &pb.MknodRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
						Pid: 1234,
					},
				},
				Name:  "test.txt",
				Mode:  0644,
				Rdev:  0,
				Umask: 022,
			},
			setupMock: func(m *mockRawFileSystem) {
				m.On("Mknod", mock.Anything, mock.MatchedBy(func(in *fuse.MknodIn) bool {
					return in.NodeId == 1 && in.Mode == 0644 && in.Rdev == 0 && in.Umask == 022
				}), "test.txt", mock.AnythingOfType("*fuse.EntryOut")).
					Return(fuse.OK)
			},
			wantErr:    false,
			wantStatus: 0,
		},
		{
			name: "not implemented",
			req: &pb.MknodRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
						Pid: 1234,
					},
				},
				Name:  "test.txt",
				Mode:  0644,
				Rdev:  0,
				Umask: 022,
			},
			setupMock: func(m *mockRawFileSystem) {
				m.On("Mknod", mock.Anything, mock.Anything, "test.txt", mock.Anything).
					Return(fuse.ENOSYS)
			},
			wantErr:     true,
			wantErrCode: codes.Unimplemented,
		},
		{
			name: "permission denied",
			req: &pb.MknodRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
						Pid: 1234,
					},
				},
				Name:  "test.txt",
				Mode:  0644,
				Rdev:  0,
				Umask: 022,
			},
			setupMock: func(m *mockRawFileSystem) {
				m.On("Mknod", mock.Anything, mock.Anything, "test.txt", mock.Anything).
					Return(fuse.EPERM)
			},
			wantErr:    false,
			wantStatus: int32(fuse.EPERM),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &mockRawFileSystem{}
			tt.setupMock(mockFS)

			s := &server{
				fs: mockFS,
			}

			resp, err := s.Mknod(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantErrCode, st.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.wantStatus, resp.Status.Code)
			}

			mockFS.AssertExpectations(t)
		})
	}
}
