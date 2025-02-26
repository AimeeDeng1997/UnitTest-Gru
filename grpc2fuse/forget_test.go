package grpc2fuse

import (
	"context"
	"testing"

	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MockRawFileSystemClient is a mock of RawFileSystemClient interface
type MockRawFileSystemClient struct {
	mock.Mock
}

func (m *MockRawFileSystemClient) Forget(ctx context.Context, in *pb.ForgetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in, opts)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

// Add other required interface methods with empty implementations
func (m *MockRawFileSystemClient) String(ctx context.Context, in *pb.StringRequest, opts ...grpc.CallOption) (*pb.StringResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Lookup(ctx context.Context, in *pb.LookupRequest, opts ...grpc.CallOption) (*pb.LookupResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) GetAttr(ctx context.Context, in *pb.GetAttrRequest, opts ...grpc.CallOption) (*pb.GetAttrResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) SetAttr(ctx context.Context, in *pb.SetAttrRequest, opts ...grpc.CallOption) (*pb.SetAttrResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Mknod(ctx context.Context, in *pb.MknodRequest, opts ...grpc.CallOption) (*pb.MknodResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Mkdir(ctx context.Context, in *pb.MkdirRequest, opts ...grpc.CallOption) (*pb.MkdirResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Unlink(ctx context.Context, in *pb.UnlinkRequest, opts ...grpc.CallOption) (*pb.UnlinkResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Rmdir(ctx context.Context, in *pb.RmdirRequest, opts ...grpc.CallOption) (*pb.RmdirResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Rename(ctx context.Context, in *pb.RenameRequest, opts ...grpc.CallOption) (*pb.RenameResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Link(ctx context.Context, in *pb.LinkRequest, opts ...grpc.CallOption) (*pb.LinkResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Symlink(ctx context.Context, in *pb.SymlinkRequest, opts ...grpc.CallOption) (*pb.SymlinkResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Readlink(ctx context.Context, in *pb.ReadlinkRequest, opts ...grpc.CallOption) (*pb.ReadlinkResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Access(ctx context.Context, in *pb.AccessRequest, opts ...grpc.CallOption) (*pb.AccessResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) GetXAttr(ctx context.Context, in *pb.GetXAttrRequest, opts ...grpc.CallOption) (*pb.GetXAttrResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) ListXAttr(ctx context.Context, in *pb.ListXAttrRequest, opts ...grpc.CallOption) (*pb.ListXAttrResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) SetXAttr(ctx context.Context, in *pb.SetXAttrRequest, opts ...grpc.CallOption) (*pb.SetXAttrResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) RemoveXAttr(ctx context.Context, in *pb.RemoveXAttrRequest, opts ...grpc.CallOption) (*pb.RemoveXAttrResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Create(ctx context.Context, in *pb.CreateRequest, opts ...grpc.CallOption) (*pb.CreateResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Open(ctx context.Context, in *pb.OpenRequest, opts ...grpc.CallOption) (*pb.OpenResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Read(ctx context.Context, in *pb.ReadRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadClient, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Lseek(ctx context.Context, in *pb.LseekRequest, opts ...grpc.CallOption) (*pb.LseekResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) GetLk(ctx context.Context, in *pb.LkRequest, opts ...grpc.CallOption) (*pb.GetLkResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) SetLk(ctx context.Context, in *pb.LkRequest, opts ...grpc.CallOption) (*pb.SetLkResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) SetLkw(ctx context.Context, in *pb.LkRequest, opts ...grpc.CallOption) (*pb.SetLkResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Release(ctx context.Context, in *pb.ReleaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Write(ctx context.Context, in *pb.WriteRequest, opts ...grpc.CallOption) (*pb.WriteResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) CopyFileRange(ctx context.Context, in *pb.CopyFileRangeRequest, opts ...grpc.CallOption) (*pb.CopyFileRangeResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Flush(ctx context.Context, in *pb.FlushRequest, opts ...grpc.CallOption) (*pb.FlushResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Fsync(ctx context.Context, in *pb.FsyncRequest, opts ...grpc.CallOption) (*pb.FsyncResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) Fallocate(ctx context.Context, in *pb.FallocateRequest, opts ...grpc.CallOption) (*pb.FallocateResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) OpenDir(ctx context.Context, in *pb.OpenDirRequest, opts ...grpc.CallOption) (*pb.OpenDirResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) ReadDir(ctx context.Context, in *pb.ReadDirRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadDirClient, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) ReadDirPlus(ctx context.Context, in *pb.ReadDirRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadDirPlusClient, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) ReleaseDir(ctx context.Context, in *pb.ReleaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) FsyncDir(ctx context.Context, in *pb.FsyncRequest, opts ...grpc.CallOption) (*pb.FsyncResponse, error) {
	return nil, nil
}
func (m *MockRawFileSystemClient) StatFs(ctx context.Context, in *pb.StatfsRequest, opts ...grpc.CallOption) (*pb.StatfsResponse, error) {
	return nil, nil
}

func TestFileSystem_Forget(t *testing.T) {
	mockClient := new(MockRawFileSystemClient)
	fs := &fileSystem{
		client: mockClient,
		opts:   []grpc.CallOption{},
	}

	tests := []struct {
		name    string
		nodeid  uint64
		nlookup uint64
	}{
		{
			name:    "basic forget",
			nodeid:  123,
			nlookup: 1,
		},
		{
			name:    "zero nodeid",
			nodeid:  0,
			nlookup: 1,
		},
		{
			name:    "multiple lookups",
			nodeid:  456,
			nlookup: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedRequest := &pb.ForgetRequest{
				Nodeid:  tt.nodeid,
				Nlookup: tt.nlookup,
			}

			mockClient.On("Forget", mock.Anything, expectedRequest, mock.Anything).
				Return(&emptypb.Empty{}, nil).Once()

			fs.Forget(tt.nodeid, tt.nlookup)

			mockClient.AssertExpectations(t)
		})
	}
}
