package grpc2fuse

import (
	"context"
	"testing"

	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type mockRawFileSystemClient struct {
	mock.Mock
}

func (m *mockRawFileSystemClient) String(ctx context.Context, in *pb.StringRequest, opts ...grpc.CallOption) (*pb.StringResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.StringResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Lookup(ctx context.Context, in *pb.LookupRequest, opts ...grpc.CallOption) (*pb.LookupResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.LookupResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Forget(ctx context.Context, in *pb.ForgetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*emptypb.Empty), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) GetAttr(ctx context.Context, in *pb.GetAttrRequest, opts ...grpc.CallOption) (*pb.GetAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.GetAttrResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) SetAttr(ctx context.Context, in *pb.SetAttrRequest, opts ...grpc.CallOption) (*pb.SetAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.SetAttrResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Mknod(ctx context.Context, in *pb.MknodRequest, opts ...grpc.CallOption) (*pb.MknodResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.MknodResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Mkdir(ctx context.Context, in *pb.MkdirRequest, opts ...grpc.CallOption) (*pb.MkdirResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.MkdirResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Unlink(ctx context.Context, in *pb.UnlinkRequest, opts ...grpc.CallOption) (*pb.UnlinkResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.UnlinkResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Rmdir(ctx context.Context, in *pb.RmdirRequest, opts ...grpc.CallOption) (*pb.RmdirResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.RmdirResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Rename(ctx context.Context, in *pb.RenameRequest, opts ...grpc.CallOption) (*pb.RenameResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.RenameResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Link(ctx context.Context, in *pb.LinkRequest, opts ...grpc.CallOption) (*pb.LinkResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.LinkResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Symlink(ctx context.Context, in *pb.SymlinkRequest, opts ...grpc.CallOption) (*pb.SymlinkResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.SymlinkResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Readlink(ctx context.Context, in *pb.ReadlinkRequest, opts ...grpc.CallOption) (*pb.ReadlinkResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.ReadlinkResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Access(ctx context.Context, in *pb.AccessRequest, opts ...grpc.CallOption) (*pb.AccessResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.AccessResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) GetXAttr(ctx context.Context, in *pb.GetXAttrRequest, opts ...grpc.CallOption) (*pb.GetXAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.GetXAttrResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) ListXAttr(ctx context.Context, in *pb.ListXAttrRequest, opts ...grpc.CallOption) (*pb.ListXAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.ListXAttrResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) SetXAttr(ctx context.Context, in *pb.SetXAttrRequest, opts ...grpc.CallOption) (*pb.SetXAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.SetXAttrResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) RemoveXAttr(ctx context.Context, in *pb.RemoveXAttrRequest, opts ...grpc.CallOption) (*pb.RemoveXAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.RemoveXAttrResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Create(ctx context.Context, in *pb.CreateRequest, opts ...grpc.CallOption) (*pb.CreateResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.CreateResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Open(ctx context.Context, in *pb.OpenRequest, opts ...grpc.CallOption) (*pb.OpenResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.OpenResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Read(ctx context.Context, in *pb.ReadRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadClient, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(pb.RawFileSystem_ReadClient), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Write(ctx context.Context, in *pb.WriteRequest, opts ...grpc.CallOption) (*pb.WriteResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.WriteResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Lseek(ctx context.Context, in *pb.LseekRequest, opts ...grpc.CallOption) (*pb.LseekResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.LseekResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Release(ctx context.Context, in *pb.ReleaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*emptypb.Empty), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) GetLk(ctx context.Context, in *pb.LkRequest, opts ...grpc.CallOption) (*pb.GetLkResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.GetLkResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) SetLk(ctx context.Context, in *pb.LkRequest, opts ...grpc.CallOption) (*pb.SetLkResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.SetLkResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) SetLkw(ctx context.Context, in *pb.LkRequest, opts ...grpc.CallOption) (*pb.SetLkResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.SetLkResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) CopyFileRange(ctx context.Context, in *pb.CopyFileRangeRequest, opts ...grpc.CallOption) (*pb.CopyFileRangeResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.CopyFileRangeResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Flush(ctx context.Context, in *pb.FlushRequest, opts ...grpc.CallOption) (*pb.FlushResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.FlushResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Fsync(ctx context.Context, in *pb.FsyncRequest, opts ...grpc.CallOption) (*pb.FsyncResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.FsyncResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) Fallocate(ctx context.Context, in *pb.FallocateRequest, opts ...grpc.CallOption) (*pb.FallocateResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.FallocateResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) OpenDir(ctx context.Context, in *pb.OpenDirRequest, opts ...grpc.CallOption) (*pb.OpenDirResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.OpenDirResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) ReadDir(ctx context.Context, in *pb.ReadDirRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadDirClient, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(pb.RawFileSystem_ReadDirClient), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) ReadDirPlus(ctx context.Context, in *pb.ReadDirRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadDirPlusClient, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(pb.RawFileSystem_ReadDirPlusClient), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) ReleaseDir(ctx context.Context, in *pb.ReleaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*emptypb.Empty), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) FsyncDir(ctx context.Context, in *pb.FsyncRequest, opts ...grpc.CallOption) (*pb.FsyncResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.FsyncResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRawFileSystemClient) StatFs(ctx context.Context, in *pb.StatfsRequest, opts ...grpc.CallOption) (*pb.StatfsResponse, error) {
	args := m.Called(ctx, in, opts)
	if resp := args.Get(0); resp != nil {
		return resp.(*pb.StatfsResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestMknod(t *testing.T) {
	tests := []struct {
		name     string
		input    *fuse.MknodIn
		fileName string
		mockResp *pb.MknodResponse
		mockErr  error
		want     fuse.Status
	}{
		{
			name: "successful mknod",
			input: &fuse.MknodIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Mode:  0644,
				Rdev:  0,
				Umask: 0022,
			},
			fileName: "testfile",
			mockResp: &pb.MknodResponse{
				Status: &pb.Status{
					Code: 0,
				},
				EntryOut: &pb.EntryOut{
					NodeId: 2,
					Attr: &pb.Attr{
						Mode:  0644,
						Owner: &pb.Owner{},
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
				Mode:  0644,
				Rdev:  0,
				Umask: 0022,
			},
			fileName: "testfile",
			mockResp: &pb.MknodResponse{
				Status: &pb.Status{
					Code: int32(fuse.EPERM),
				},
				EntryOut: &pb.EntryOut{
					Attr: &pb.Attr{
						Owner: &pb.Owner{},
					},
				},
			},
			mockErr: nil,
			want:    fuse.EPERM,
		},
		{
			name: "grpc error",
			input: &fuse.MknodIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Mode:  0644,
				Rdev:  0,
				Umask: 0022,
			},
			fileName: "testfile",
			mockResp: nil,
			mockErr:  grpc.ErrServerStopped,
			want:     fuse.EIO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(mockRawFileSystemClient)
			fs := &fileSystem{
				client: mockClient,
			}

			expectedRequest := &pb.MknodRequest{
				Header: toPbHeader(&tt.input.InHeader),
				Name:   tt.fileName,
				Mode:   tt.input.Mode,
				Rdev:   tt.input.Rdev,
				Umask:  tt.input.Umask,
			}

			mockClient.On("Mknod", mock.Anything, expectedRequest, mock.Anything).Return(tt.mockResp, tt.mockErr)

			var out fuse.EntryOut
			got := fs.Mknod(make(<-chan struct{}), tt.input, tt.fileName, &out)

			assert.Equal(t, tt.want, got)
			if tt.mockResp != nil && tt.mockResp.Status.Code == 0 {
				assert.Equal(t, tt.mockResp.EntryOut.NodeId, out.NodeId)
				assert.Equal(t, tt.mockResp.EntryOut.Attr.Mode, out.Attr.Mode)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
