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

// Mock RawFileSystemClient
type mockRawFileSystemClient struct {
	mock.Mock
	pb.RawFileSystemClient
}

func (m *mockRawFileSystemClient) SetXAttr(ctx context.Context, in *pb.SetXAttrRequest, opts ...grpc.CallOption) (*pb.SetXAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.SetXAttrResponse), args.Error(1)
}

func TestFileSystem_SetXAttr(t *testing.T) {
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
				Size:  10,
				Flags: 0,
			},
			attr: "user.test",
			data: []byte("test_value"),
			mockResp: &pb.SetXAttrResponse{
				Status: &pb.Status{
					Code: 0,
				},
			},
			mockErr: nil,
			want:    fuse.OK,
		},
		{
			name: "error from server",
			input: &fuse.SetXAttrIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Size:  10,
				Flags: 0,
			},
			attr: "user.test",
			data: []byte("test_value"),
			mockResp: &pb.SetXAttrResponse{
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
			mockClient := new(mockRawFileSystemClient)

			fs := &fileSystem{
				client: mockClient,
			}

			// Setup expectations
			mockClient.On("SetXAttr", mock.Anything, &pb.SetXAttrRequest{
				Header: toPbHeader(&tt.input.InHeader),
				Attr:   tt.attr,
				Data:   tt.data,
				Size:   tt.input.Size,
				Flags:  tt.input.Flags,
			}, mock.Anything).Return(tt.mockResp, tt.mockErr)

			// Create cancel channel
			cancel := make(chan struct{})

			// Call the method
			got := fs.SetXAttr(cancel, tt.input, tt.attr, tt.data)

			// Assert result
			assert.Equal(t, tt.want, got)

			// Verify expectations
			mockClient.AssertExpectations(t)
		})
	}
}
