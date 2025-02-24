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

type MockRawFileSystem struct {
	mock.Mock
}

func (m *MockRawFileSystem) Write(cancel <-chan struct{}, input *fuse.WriteIn, data []byte) (written uint32, code fuse.Status) {
	args := m.Called(cancel, input, data)
	return args.Get(0).(uint32), args.Get(1).(fuse.Status)
}

func (m *MockRawFileSystem) String() string {
	return "MockRawFileSystem"
}

func TestWrite(t *testing.T) {
	mockFS := &MockRawFileSystem{}
	server := fuse2grpc.NewServer(mockFS)

	tests := []struct {
		name           string
		req           *pb.WriteRequest
		mockWritten   uint32
		mockStatus    fuse.Status
		expectedError error
		expectedResp  *pb.WriteResponse
	}{
		{
			name: "successful write",
			req: &pb.WriteRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Fh:         2,
				Offset:     100,
				Size:       10,
				WriteFlags: 0,
				Data:       []byte("test data"),
			},
			mockWritten: 10,
			mockStatus:  fuse.OK,
			expectedResp: &pb.WriteResponse{
				Written: 10,
				Status:  &pb.Status{Code: 0},
			},
		},
		{
			name: "write not implemented",
			req: &pb.WriteRequest{
				Header: &pb.InHeader{NodeId: 1},
				Data:   []byte("test"),
			},
			mockStatus:    fuse.ENOSYS,
			expectedError: status.Errorf(codes.Unimplemented, "method Write not implemented"),
		},
		{
			name: "write error",
			req: &pb.WriteRequest{
				Header: &pb.InHeader{NodeId: 1},
				Data:   []byte("test"),
			},
			mockWritten: 0,
			mockStatus:  fuse.EIO,
			expectedResp: &pb.WriteResponse{
				Written: 0,
				Status:  &pb.Status{Code: int32(fuse.EIO)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS.On("Write", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockWritten, tt.mockStatus).Once()

			resp, err := server.Write(context.Background(), tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}

			mockFS.AssertExpectations(t)
		})
	}
}
