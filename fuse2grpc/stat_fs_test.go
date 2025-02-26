package fuse2grpc_test

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/fuse2grpc"
	"github.com/chiyutianyi/grpcfuse/pb"
)

type mockRawFileSystem struct {
	status fuse.Status
}

func (m *mockRawFileSystem) String() string { return "mock" }
func (m *mockRawFileSystem) SetDebug(debug bool) {}
func (m *mockRawFileSystem) Init(server *fuse.Server) {}
func (m *mockRawFileSystem) Lookup(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Forget(nodeID, nlookup uint64) {}
func (m *mockRawFileSystem) GetAttr(cancel <-chan struct{}, input *fuse.GetAttrIn, out *fuse.AttrOut) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) SetAttr(cancel <-chan struct{}, input *fuse.SetAttrIn, out *fuse.AttrOut) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Mknod(cancel <-chan struct{}, input *fuse.MknodIn, name string, out *fuse.EntryOut) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Mkdir(cancel <-chan struct{}, input *fuse.MkdirIn, name string, out *fuse.EntryOut) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Unlink(cancel <-chan struct{}, header *fuse.InHeader, name string) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Rmdir(cancel <-chan struct{}, header *fuse.InHeader, name string) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Rename(cancel <-chan struct{}, input *fuse.RenameIn, oldName string, newName string) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Link(cancel <-chan struct{}, input *fuse.LinkIn, name string, out *fuse.EntryOut) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Symlink(cancel <-chan struct{}, header *fuse.InHeader, pointedTo string, linkName string, out *fuse.EntryOut) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Readlink(cancel <-chan struct{}, header *fuse.InHeader) ([]byte, fuse.Status) {
	return nil, fuse.OK
}
func (m *mockRawFileSystem) Access(cancel <-chan struct{}, input *fuse.AccessIn) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) GetXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string, dest []byte) (uint32, fuse.Status) {
	return 0, fuse.OK
}
func (m *mockRawFileSystem) ListXAttr(cancel <-chan struct{}, header *fuse.InHeader, dest []byte) (uint32, fuse.Status) {
	return 0, fuse.OK
}
func (m *mockRawFileSystem) SetXAttr(cancel <-chan struct{}, input *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) RemoveXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Create(cancel <-chan struct{}, input *fuse.CreateIn, name string, out *fuse.CreateOut) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Open(cancel <-chan struct{}, input *fuse.OpenIn, out *fuse.OpenOut) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Read(cancel <-chan struct{}, input *fuse.ReadIn, buf []byte) (fuse.ReadResult, fuse.Status) {
	return nil, fuse.OK
}
func (m *mockRawFileSystem) Lseek(cancel <-chan struct{}, in *fuse.LseekIn, out *fuse.LseekOut) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) GetLk(cancel <-chan struct{}, input *fuse.LkIn, out *fuse.LkOut) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) SetLk(cancel <-chan struct{}, input *fuse.LkIn) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) SetLkw(cancel <-chan struct{}, input *fuse.LkIn) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Release(cancel <-chan struct{}, input *fuse.ReleaseIn) {}
func (m *mockRawFileSystem) Write(cancel <-chan struct{}, input *fuse.WriteIn, data []byte) (uint32, fuse.Status) {
	return 0, fuse.OK
}
func (m *mockRawFileSystem) CopyFileRange(cancel <-chan struct{}, input *fuse.CopyFileRangeIn) (uint32, fuse.Status) {
	return 0, fuse.OK
}
func (m *mockRawFileSystem) Flush(cancel <-chan struct{}, input *fuse.FlushIn) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Fsync(cancel <-chan struct{}, input *fuse.FsyncIn) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) Fallocate(cancel <-chan struct{}, input *fuse.FallocateIn) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) OpenDir(cancel <-chan struct{}, input *fuse.OpenIn, out *fuse.OpenOut) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) ReadDir(cancel <-chan struct{}, input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) ReadDirPlus(cancel <-chan struct{}, input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) ReleaseDir(input *fuse.ReleaseIn) {}
func (m *mockRawFileSystem) FsyncDir(cancel <-chan struct{}, input *fuse.FsyncIn) fuse.Status {
	return fuse.OK
}
func (m *mockRawFileSystem) StatFs(cancel <-chan struct{}, header *fuse.InHeader, out *fuse.StatfsOut) fuse.Status {
	if m.status == fuse.OK {
		out.Blocks = 1000
		out.Bfree = 500
		out.Bavail = 400
		out.Files = 100
		out.Ffree = 50
		out.Bsize = 4096
		out.NameLen = 255
		out.Frsize = 4096
	}
	return m.status
}

func (m *mockRawFileSystem) Destroy() {}

func TestStatFs(t *testing.T) {
	tests := []struct {
		name           string
		status         fuse.Status
		expectedError  error
		expectedStatus *pb.Status
		expectedResp   *pb.StatfsResponse
	}{
		{
			name:          "success",
			status:        fuse.OK,
			expectedError: nil,
			expectedResp: &pb.StatfsResponse{
				Status:  &pb.Status{Code: 0},
				Blocks:  1000,
				Bfree:   500,
				Bavail:  400,
				Files:   100,
				Ffree:   50,
				Bsize:   4096,
				NameLen: 255,
				Frsize:  4096,
			},
		},
		{
			name:           "not implemented",
			status:         fuse.ENOSYS,
			expectedError:  status.Errorf(codes.Unimplemented, "method StatFS not implemented"),
			expectedStatus: nil,
		},
		{
			name:   "error status",
			status: fuse.EACCES,
			expectedResp: &pb.StatfsResponse{
				Status: &pb.Status{Code: int32(fuse.EACCES)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRawFileSystem{status: tt.status}
			server := fuse2grpc.NewServer(mock)

			req := &pb.StatfsRequest{
				Input: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
						Pid: 1234,
					},
				},
			}

			resp, err := server.StatFs(context.Background(), req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, resp)
			}
		})
	}
}
