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

type MockRawFileSystemClient struct {
	mock.Mock
}

// Implement all required methods of pb.RawFileSystemClient interface
func (m *MockRawFileSystemClient) String(ctx context.Context, in *pb.StringRequest, opts ...grpc.CallOption) (*pb.StringResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.StringResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Lookup(ctx context.Context, in *pb.LookupRequest, opts ...grpc.CallOption) (*pb.LookupResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.LookupResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Forget(ctx context.Context, in *pb.ForgetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (m *MockRawFileSystemClient) GetAttr(ctx context.Context, in *pb.GetAttrRequest, opts ...grpc.CallOption) (*pb.GetAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.GetAttrResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) SetAttr(ctx context.Context, in *pb.SetAttrRequest, opts ...grpc.CallOption) (*pb.SetAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.SetAttrResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Mknod(ctx context.Context, in *pb.MknodRequest, opts ...grpc.CallOption) (*pb.MknodResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.MknodResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Mkdir(ctx context.Context, in *pb.MkdirRequest, opts ...grpc.CallOption) (*pb.MkdirResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.MkdirResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Unlink(ctx context.Context, in *pb.UnlinkRequest, opts ...grpc.CallOption) (*pb.UnlinkResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.UnlinkResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Rmdir(ctx context.Context, in *pb.RmdirRequest, opts ...grpc.CallOption) (*pb.RmdirResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.RmdirResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Rename(ctx context.Context, in *pb.RenameRequest, opts ...grpc.CallOption) (*pb.RenameResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.RenameResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Link(ctx context.Context, in *pb.LinkRequest, opts ...grpc.CallOption) (*pb.LinkResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.LinkResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Symlink(ctx context.Context, in *pb.SymlinkRequest, opts ...grpc.CallOption) (*pb.SymlinkResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.SymlinkResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Readlink(ctx context.Context, in *pb.ReadlinkRequest, opts ...grpc.CallOption) (*pb.ReadlinkResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.ReadlinkResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Access(ctx context.Context, in *pb.AccessRequest, opts ...grpc.CallOption) (*pb.AccessResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.AccessResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) GetXAttr(ctx context.Context, in *pb.GetXAttrRequest, opts ...grpc.CallOption) (*pb.GetXAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.GetXAttrResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) ListXAttr(ctx context.Context, in *pb.ListXAttrRequest, opts ...grpc.CallOption) (*pb.ListXAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.ListXAttrResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) SetXAttr(ctx context.Context, in *pb.SetXAttrRequest, opts ...grpc.CallOption) (*pb.SetXAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.SetXAttrResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) RemoveXAttr(ctx context.Context, in *pb.RemoveXAttrRequest, opts ...grpc.CallOption) (*pb.RemoveXAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.RemoveXAttrResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Create(ctx context.Context, in *pb.CreateRequest, opts ...grpc.CallOption) (*pb.CreateResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.CreateResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Open(ctx context.Context, in *pb.OpenRequest, opts ...grpc.CallOption) (*pb.OpenResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.OpenResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Read(ctx context.Context, in *pb.ReadRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadClient, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(pb.RawFileSystem_ReadClient), args.Error(1)
}

func (m *MockRawFileSystemClient) Write(ctx context.Context, in *pb.WriteRequest, opts ...grpc.CallOption) (*pb.WriteResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.WriteResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Lseek(ctx context.Context, in *pb.LseekRequest, opts ...grpc.CallOption) (*pb.LseekResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.LseekResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Release(ctx context.Context, in *pb.ReleaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (m *MockRawFileSystemClient) GetLk(ctx context.Context, in *pb.LkRequest, opts ...grpc.CallOption) (*pb.GetLkResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.GetLkResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) SetLk(ctx context.Context, in *pb.LkRequest, opts ...grpc.CallOption) (*pb.SetLkResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.SetLkResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) SetLkw(ctx context.Context, in *pb.LkRequest, opts ...grpc.CallOption) (*pb.SetLkResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.SetLkResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) CopyFileRange(ctx context.Context, in *pb.CopyFileRangeRequest, opts ...grpc.CallOption) (*pb.CopyFileRangeResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.CopyFileRangeResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Flush(ctx context.Context, in *pb.FlushRequest, opts ...grpc.CallOption) (*pb.FlushResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.FlushResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Fsync(ctx context.Context, in *pb.FsyncRequest, opts ...grpc.CallOption) (*pb.FsyncResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.FsyncResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Fallocate(ctx context.Context, in *pb.FallocateRequest, opts ...grpc.CallOption) (*pb.FallocateResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.FallocateResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) OpenDir(ctx context.Context, in *pb.OpenDirRequest, opts ...grpc.CallOption) (*pb.OpenDirResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.OpenDirResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) ReadDir(ctx context.Context, in *pb.ReadDirRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadDirClient, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(pb.RawFileSystem_ReadDirClient), args.Error(1)
}

func (m *MockRawFileSystemClient) ReadDirPlus(ctx context.Context, in *pb.ReadDirRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadDirPlusClient, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(pb.RawFileSystem_ReadDirPlusClient), args.Error(1)
}

func (m *MockRawFileSystemClient) ReleaseDir(ctx context.Context, in *pb.ReleaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (m *MockRawFileSystemClient) FsyncDir(ctx context.Context, in *pb.FsyncRequest, opts ...grpc.CallOption) (*pb.FsyncResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.FsyncResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) StatFs(ctx context.Context, in *pb.StatfsRequest, opts ...grpc.CallOption) (*pb.StatfsResponse, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*pb.StatfsResponse), args.Error(1)
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
			name: "error setting xattr",
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
			mockClient := new(MockRawFileSystemClient)
			fs := &fileSystem{
				client: mockClient,
			}

			expectedRequest := &pb.SetXAttrRequest{
				Header: toPbHeader(&tt.input.InHeader),
				Attr:   tt.attr,
				Data:   tt.data,
				Size:   tt.input.Size,
				Flags:  tt.input.Flags,
			}

			mockClient.On("SetXAttr", mock.Anything, expectedRequest, mock.Anything).Return(tt.mockResp, tt.mockErr)

			got := fs.SetXAttr(make(chan struct{}), tt.input, tt.attr, tt.data)
			assert.Equal(t, tt.want, got)

			mockClient.AssertExpectations(t)
		})
	}
}
