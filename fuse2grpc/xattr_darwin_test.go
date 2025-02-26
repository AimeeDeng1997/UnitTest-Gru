package fuse2grpc_test

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/fuse2grpc"
	"github.com/chiyutianyi/grpcfuse/pb"
)

type mockFS struct {
	mock.Mock
}

func (m *mockFS) SetXAttr(cancel <-chan struct{}, in *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
	args := m.Called(cancel, in, attr, data)
	return args.Get(0).(fuse.Status)
}

func (m *mockFS) String() string {
	return "mockFS"
}

func TestSetXAttr(t *testing.T) {
	tests := []struct {
		name        string
		req         *pb.SetXAttrRequest
		setupMock   func(*mockFS)
		wantStatus  *pb.Status
		wantErrCode codes.Code
	}{
		{
			name: "successful set xattr",
			req: &pb.SetXAttrRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Attr:     "user.test",
				Data:     []byte("test_value"),
				Size:     10,
				Flags:    0,
				Position: 0,
				Padding:  0,
			},
			setupMock: func(m *mockFS) {
				m.On("SetXAttr", mock.Anything, mock.MatchedBy(func(in *fuse.SetXAttrIn) bool {
					return in.NodeId == 1
				}), "user.test", []byte("test_value")).Return(fuse.OK)
			},
			wantStatus: &pb.Status{Code: 0},
		},
		{
			name: "not implemented",
			req: &pb.SetXAttrRequest{
				Header: &pb.InHeader{NodeId: 1},
				Attr:   "user.test",
			},
			setupMock: func(m *mockFS) {
				m.On("SetXAttr", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fuse.ENOSYS)
			},
			wantErrCode: codes.Unimplemented,
		},
		{
			name: "permission denied",
			req: &pb.SetXAttrRequest{
				Header: &pb.InHeader{NodeId: 1},
				Attr:   "user.test",
			},
			setupMock: func(m *mockFS) {
				m.On("SetXAttr", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fuse.EPERM)
			},
			wantStatus: &pb.Status{Code: int32(fuse.EPERM)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFs := &mockFS{}
			tt.setupMock(mockFs)

			server := fuse2grpc.NewServer(mockFs)
			resp, err := server.SetXAttr(context.Background(), tt.req)

			if tt.wantErrCode != 0 {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantErrCode, st.Code())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, resp.Status)
			}

			mockFs.AssertExpectations(t)
		})
	}
}
