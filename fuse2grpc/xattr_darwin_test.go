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
		input       *pb.SetXAttrRequest
		setupMock   func(*mockFS)
		expectError error
		expectResp  *pb.SetXAttrResponse
	}{
		{
			name: "successful set xattr",
			input: &pb.SetXAttrRequest{
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
				m.On("SetXAttr",
					mock.Anything,
					mock.MatchedBy(func(in *fuse.SetXAttrIn) bool {
						return in.NodeId == 1 && in.Size == 10
					}),
					"user.test",
					[]byte("test_value"),
				).Return(fuse.OK)
			},
			expectResp: &pb.SetXAttrResponse{
				Status: &pb.Status{Code: 0},
			},
		},
		{
			name: "not implemented",
			input: &pb.SetXAttrRequest{
				Header: &pb.InHeader{NodeId: 1},
				Attr:   "user.test",
			},
			setupMock: func(m *mockFS) {
				m.On("SetXAttr",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(fuse.ENOSYS)
			},
			expectError: status.Errorf(codes.Unimplemented, "method GetXAttr not implemented"),
		},
		{
			name: "operation failed",
			input: &pb.SetXAttrRequest{
				Header: &pb.InHeader{NodeId: 1},
				Attr:   "user.test",
			},
			setupMock: func(m *mockFS) {
				m.On("SetXAttr",
					mock.Anything,
					mock.Anything,
					mock.Anything,
					mock.Anything,
				).Return(fuse.EACCES)
			},
			expectResp: &pb.SetXAttrResponse{
				Status: &pb.Status{Code: int32(fuse.EACCES)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFs := &mockFS{}
			tt.setupMock(mockFs)

			server := fuse2grpc.NewServer(mockFs)
			resp, err := server.SetXAttr(context.Background(), tt.input)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectError.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectResp, resp)
			}

			mockFs.AssertExpectations(t)
		})
	}
}
