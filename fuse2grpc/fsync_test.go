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

func (m *mockFS) Init(server *fuse.Server) {}

func (m *mockFS) Fsync(cancel <-chan struct{}, in *fuse.FsyncIn) fuse.Status {
	args := m.Called(cancel, in)
	return args.Get(0).(fuse.Status)
}

func (m *mockFS) String() string                                    { return "mockFS" }
func (m *mockFS) SetDebug(debug bool)                              {}
func (m *mockFS) Lookup(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Forget(nodeID uint64, nlookup uint64)            {}
func (m *mockFS) GetAttr(cancel <-chan struct{}, input *fuse.GetAttrIn, out *fuse.AttrOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) SetAttr(cancel <-chan struct{}, input *fuse.SetAttrIn, out *fuse.AttrOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Mknod(cancel <-chan struct{}, input *fuse.MknodIn, name string, out *fuse.EntryOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Mkdir(cancel <-chan struct{}, input *fuse.MkdirIn, name string, out *fuse.EntryOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Unlink(cancel <-chan struct{}, header *fuse.InHeader, name string) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Rmdir(cancel <-chan struct{}, header *fuse.InHeader, name string) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Rename(cancel <-chan struct{}, input *fuse.RenameIn, oldName string, newName string) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Link(cancel <-chan struct{}, input *fuse.LinkIn, name string, out *fuse.EntryOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Symlink(cancel <-chan struct{}, header *fuse.InHeader, pointedTo string, linkName string, out *fuse.EntryOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Readlink(cancel <-chan struct{}, header *fuse.InHeader) ([]byte, fuse.Status) {
	return nil, fuse.OK
}
func (m *mockFS) Access(cancel <-chan struct{}, input *fuse.AccessIn) fuse.Status {
	return fuse.OK
}
func (m *mockFS) GetXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string, dest []byte) (uint32, fuse.Status) {
	return 0, fuse.OK
}
func (m *mockFS) ListXAttr(cancel <-chan struct{}, header *fuse.InHeader, dest []byte) (uint32, fuse.Status) {
	return 0, fuse.OK
}
func (m *mockFS) SetXAttr(cancel <-chan struct{}, input *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
	return fuse.OK
}
func (m *mockFS) RemoveXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Create(cancel <-chan struct{}, input *fuse.CreateIn, name string, out *fuse.CreateOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Open(cancel <-chan struct{}, input *fuse.OpenIn, out *fuse.OpenOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Read(cancel <-chan struct{}, input *fuse.ReadIn, buf []byte) (fuse.ReadResult, fuse.Status) {
	return nil, fuse.OK
}
func (m *mockFS) Lseek(cancel <-chan struct{}, in *fuse.LseekIn, out *fuse.LseekOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) GetLk(cancel <-chan struct{}, input *fuse.LkIn, out *fuse.LkOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) SetLk(cancel <-chan struct{}, input *fuse.LkIn) fuse.Status {
	return fuse.OK
}
func (m *mockFS) SetLkw(cancel <-chan struct{}, input *fuse.LkIn) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Release(cancel <-chan struct{}, input *fuse.ReleaseIn)     {}
func (m *mockFS) Write(cancel <-chan struct{}, input *fuse.WriteIn, data []byte) (uint32, fuse.Status) {
	return 0, fuse.OK
}
func (m *mockFS) CopyFileRange(cancel <-chan struct{}, input *fuse.CopyFileRangeIn) (uint32, fuse.Status) {
	return 0, fuse.OK
}
func (m *mockFS) Flush(cancel <-chan struct{}, input *fuse.FlushIn) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Fallocate(cancel <-chan struct{}, input *fuse.FallocateIn) fuse.Status {
	return fuse.OK
}
func (m *mockFS) OpenDir(cancel <-chan struct{}, input *fuse.OpenIn, out *fuse.OpenOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) ReadDir(cancel <-chan struct{}, input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	return fuse.OK
}
func (m *mockFS) ReadDirPlus(cancel <-chan struct{}, input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	return fuse.OK
}
func (m *mockFS) ReleaseDir(input *fuse.ReleaseIn) {}
func (m *mockFS) FsyncDir(cancel <-chan struct{}, input *fuse.FsyncIn) fuse.Status {
	return fuse.OK
}
func (m *mockFS) StatFs(cancel <-chan struct{}, input *fuse.InHeader, out *fuse.StatfsOut) fuse.Status {
	return fuse.OK
}

func TestFsync(t *testing.T) {
	mockfs := &mockFS{}
	server := fuse2grpc.NewServer(mockfs)

	tests := []struct {
		name        string
		req         *pb.FsyncRequest
		mockStatus  fuse.Status
		wantStatus  int32
		wantErrCode codes.Code
	}{
		{
			name: "successful fsync",
			req: &pb.FsyncRequest{
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
				FsyncFlags: 3,
				Padding:    4,
			},
			mockStatus:  fuse.OK,
			wantStatus:  0,
			wantErrCode: codes.OK,
		},
		{
			name: "unimplemented fsync",
			req: &pb.FsyncRequest{
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
				FsyncFlags: 3,
				Padding:    4,
			},
			mockStatus:  fuse.ENOSYS,
			wantStatus:  0,
			wantErrCode: codes.Unimplemented,
		},
		{
			name: "failed fsync",
			req: &pb.FsyncRequest{
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
				FsyncFlags: 3,
				Padding:    4,
			},
			mockStatus:  fuse.EINVAL,
			wantStatus:  int32(fuse.EINVAL),
			wantErrCode: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockfs.On("Fsync", mock.Anything, mock.MatchedBy(func(in *fuse.FsyncIn) bool {
				return in.NodeId == tt.req.Header.NodeId &&
					in.Fh == tt.req.Fh &&
					in.FsyncFlags == tt.req.FsyncFlags &&
					in.Padding == tt.req.Padding
			})).Return(tt.mockStatus).Once()

			resp, err := server.Fsync(context.Background(), tt.req)

			if tt.wantErrCode != codes.OK {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantErrCode, st.Code())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStatus, resp.Status.Code)
			}

			mockfs.AssertExpectations(t)
		})
	}
}
