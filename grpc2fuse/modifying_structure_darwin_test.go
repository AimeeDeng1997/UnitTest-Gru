package grpc2fuse_test

import (
	"context"
	"testing"

	"github.com/chiyutianyi/grpcfuse/grpc2fuse"
	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type mockRawFileSystemClient struct {
	mock.Mock
}

func (m *mockRawFileSystemClient) Mknod(ctx context.Context, in *pb.MknodRequest, opts ...grpc.CallOption) (*pb.MknodResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.MknodResponse), args.Error(1)
}

// Implement other required interface methods with empty implementations
func (m *mockRawFileSystemClient) String(ctx context.Context, in *pb.StringRequest, opts ...grpc.CallOption) (*pb.StringResponse, error) {
	return nil, nil
}

func (m *mockRawFileSystemClient) Lookup(ctx context.Context, in *pb.LookupRequest, opts ...grpc.CallOption) (*pb.LookupResponse, error) {
	return nil, nil
}

func (m *mockRawFileSystemClient) Forget(ctx context.Context, in *pb.ForgetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}

func (m *mockRawFileSystemClient) GetAttr(ctx context.Context, in *pb.GetAttrRequest, opts ...grpc.CallOption) (*pb.GetAttrResponse, error) {
	return nil, nil
}

func TestMknod(t *testing.T) {
	tests := []struct {
		name     string
		input    *fuse.MknodIn
		nodeName string
		mockResp *pb.MknodResponse
		mockErr  error
		want     fuse.Status
	}{
		{
			name: "successful mknod",
			input: &fuse.MknodIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
					Uid:    1000,
					Gid:    1000,
					Pid:    12345,
				},
				Mode: 0644,
				Rdev: 0,
			},
			nodeName: "testnode",
			mockResp: &pb.MknodResponse{
				Status: &pb.Status{
					Code: 0,
				},
				EntryOut: &pb.EntryOut{
					NodeId:          2,
					Generation:      1,
					EntryValid:     3600,
					AttrValid:      3600,
					EntryValidNsec: 0,
					AttrValidNsec:  0,
					Attr: &pb.Attr{
						Ino:   2,
						Mode:  0644,
						Nlink: 1,
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
			},
			mockErr: nil,
			want:    fuse.OK,
		},
		{
			name: "error response from server",
			input: &fuse.MknodIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Mode: 0644,
				Rdev: 0,
			},
			nodeName: "testnode",
			mockResp: &pb.MknodResponse{
				Status: &pb.Status{
					Code: int32(fuse.EPERM),
				},
			},
			mockErr: nil,
			want:    fuse.EPERM,
		},
		{
			name: "invalid mode",
			input: &fuse.MknodIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Mode: 0777,
				Rdev: 0,
			},
			nodeName: "testnode",
			mockResp: &pb.MknodResponse{
				Status: &pb.Status{
					Code: int32(fuse.EINVAL),
				},
			},
			mockErr: nil,
			want:    fuse.EINVAL,
		},
		{
			name: "node already exists",
			input: &fuse.MknodIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Mode: 0644,
				Rdev: 0,
			},
			nodeName: "existingnode",
			mockResp: &pb.MknodResponse{
				Status: &pb.Status{
					Code: int32(fuse.EEXIST),
				},
			},
			mockErr: nil,
			want:    fuse.EEXIST,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockRawFileSystemClient{}
			fs := grpc2fuse.NewFileSystem(mockClient)

			expectedRequest := &pb.MknodRequest{
				Header: &pb.InHeader{
					NodeId: tt.input.NodeId,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: tt.input.Uid,
							Gid: tt.input.Gid,
						},
						Pid: tt.input.Pid,
					},
				},
				Name: tt.nodeName,
				Mode: tt.input.Mode,
				Rdev: tt.input.Rdev,
			}

			mockClient.On("Mknod", mock.Anything, expectedRequest, mock.Anything).Return(tt.mockResp, tt.mockErr)

			var out fuse.EntryOut
			cancel := make(chan struct{})
			got := fs.Mknod(cancel, tt.input, tt.nodeName, &out)

			assert.Equal(t, tt.want, got)
			if tt.want == fuse.OK {
				assert.Equal(t, uint64(tt.mockResp.EntryOut.NodeId), out.NodeId)
				assert.Equal(t, uint64(tt.mockResp.EntryOut.Generation), out.Generation)
				assert.Equal(t, uint64(tt.mockResp.EntryOut.EntryValid), out.EntryValid)
				assert.Equal(t, uint64(tt.mockResp.EntryOut.AttrValid), out.AttrValid)
				assert.Equal(t, uint32(tt.mockResp.EntryOut.EntryValidNsec), out.EntryValidNsec)
				assert.Equal(t, uint32(tt.mockResp.EntryOut.AttrValidNsec), out.AttrValidNsec)
				assert.Equal(t, uint64(tt.mockResp.EntryOut.Attr.Ino), out.Attr.Ino)
				assert.Equal(t, uint32(tt.mockResp.EntryOut.Attr.Mode), out.Attr.Mode)
				assert.Equal(t, uint32(tt.mockResp.EntryOut.Attr.Nlink), out.Attr.Nlink)
				assert.Equal(t, uint32(tt.mockResp.EntryOut.Attr.Owner.Uid), out.Attr.Uid)
				assert.Equal(t, uint32(tt.mockResp.EntryOut.Attr.Owner.Gid), out.Attr.Gid)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
