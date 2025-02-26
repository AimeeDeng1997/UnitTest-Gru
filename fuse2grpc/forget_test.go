package fuse2grpc_test

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/chiyutianyi/grpcfuse/fuse2grpc"
	"github.com/chiyutianyi/grpcfuse/pb"
)

type MockRawFileSystem struct {
	mock.Mock
}

func (m *MockRawFileSystem) Init(s *fuse.Server) {
	m.Called(s)
}

func (m *MockRawFileSystem) Forget(nodeid uint64, nlookup uint64) {
	m.Called(nodeid, nlookup)
}

func (m *MockRawFileSystem) String() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRawFileSystem) SetDebug(debug bool) {
	m.Called(debug)
}

func TestForget(t *testing.T) {
	mockFS := new(MockRawFileSystem)
	server := fuse2grpc.NewServer(mockFS)

	tests := []struct {
		name    string
		req     *pb.ForgetRequest
		setup   func()
		wantErr bool
	}{
		{
			name: "successful forget",
			req: &pb.ForgetRequest{
				Nodeid:  123,
				Nlookup: 1,
			},
			setup: func() {
				mockFS.On("Forget", uint64(123), uint64(1)).Return()
			},
			wantErr: false,
		},
		{
			name: "forget with zero nodeid",
			req: &pb.ForgetRequest{
				Nodeid:  0,
				Nlookup: 1,
			},
			setup: func() {
				mockFS.On("Forget", uint64(0), uint64(1)).Return()
			},
			wantErr: false,
		},
		{
			name: "forget with multiple lookups",
			req: &pb.ForgetRequest{
				Nodeid:  456,
				Nlookup: 5,
			},
			setup: func() {
				mockFS.On("Forget", uint64(456), uint64(5)).Return()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			_, err := server.Forget(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockFS.AssertExpectations(t)
		})
	}
}

// Minimal implementation of other required RawFileSystem interface methods
func (m *MockRawFileSystem) Lookup(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) GetAttr(cancel <-chan struct{}, input *fuse.GetAttrIn, out *fuse.AttrOut) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) SetAttr(cancel <-chan struct{}, input *fuse.SetAttrIn, out *fuse.AttrOut) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Mknod(cancel <-chan struct{}, input *fuse.MknodIn, name string, out *fuse.EntryOut) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Mkdir(cancel <-chan struct{}, input *fuse.MkdirIn, name string, out *fuse.EntryOut) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Unlink(cancel <-chan struct{}, header *fuse.InHeader, name string) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Rmdir(cancel <-chan struct{}, header *fuse.InHeader, name string) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Rename(cancel <-chan struct{}, input *fuse.RenameIn, oldName string, newName string) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Link(cancel <-chan struct{}, input *fuse.LinkIn, filename string, out *fuse.EntryOut) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Symlink(cancel <-chan struct{}, header *fuse.InHeader, pointedTo string, linkName string, out *fuse.EntryOut) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Readlink(cancel <-chan struct{}, header *fuse.InHeader) ([]byte, fuse.Status) {
	return nil, fuse.ENOSYS
}

func (m *MockRawFileSystem) Access(cancel <-chan struct{}, input *fuse.AccessIn) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) GetXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string, dest []byte) (uint32, fuse.Status) {
	return 0, fuse.ENOSYS
}

func (m *MockRawFileSystem) ListXAttr(cancel <-chan struct{}, header *fuse.InHeader, dest []byte) (uint32, fuse.Status) {
	return 0, fuse.ENOSYS
}

func (m *MockRawFileSystem) SetXAttr(cancel <-chan struct{}, input *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) RemoveXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Create(cancel <-chan struct{}, input *fuse.CreateIn, name string, out *fuse.CreateOut) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Open(cancel <-chan struct{}, input *fuse.OpenIn, out *fuse.OpenOut) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Read(cancel <-chan struct{}, input *fuse.ReadIn, buf []byte) (fuse.ReadResult, fuse.Status) {
	return nil, fuse.ENOSYS
}

func (m *MockRawFileSystem) Lseek(cancel <-chan struct{}, in *fuse.LseekIn, out *fuse.LseekOut) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) GetLk(cancel <-chan struct{}, input *fuse.LkIn, out *fuse.LkOut) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) SetLk(cancel <-chan struct{}, input *fuse.LkIn) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) SetLkw(cancel <-chan struct{}, input *fuse.LkIn) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Release(cancel <-chan struct{}, input *fuse.ReleaseIn) {
}

func (m *MockRawFileSystem) Write(cancel <-chan struct{}, input *fuse.WriteIn, data []byte) (uint32, fuse.Status) {
	return 0, fuse.ENOSYS
}

func (m *MockRawFileSystem) CopyFileRange(cancel <-chan struct{}, input *fuse.CopyFileRangeIn) (uint32, fuse.Status) {
	return 0, fuse.ENOSYS
}

func (m *MockRawFileSystem) Flush(cancel <-chan struct{}, input *fuse.FlushIn) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Fsync(cancel <-chan struct{}, input *fuse.FsyncIn) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) Fallocate(cancel <-chan struct{}, input *fuse.FallocateIn) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) OpenDir(cancel <-chan struct{}, input *fuse.OpenIn, out *fuse.OpenOut) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) ReadDir(cancel <-chan struct{}, input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) ReadDirPlus(cancel <-chan struct{}, input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) ReleaseDir(input *fuse.ReleaseIn) {
}

func (m *MockRawFileSystem) FsyncDir(cancel <-chan struct{}, input *fuse.FsyncIn) fuse.Status {
	return fuse.ENOSYS
}

func (m *MockRawFileSystem) StatFs(cancel <-chan struct{}, input *fuse.InHeader, out *fuse.StatfsOut) fuse.Status {
	return fuse.ENOSYS
}
