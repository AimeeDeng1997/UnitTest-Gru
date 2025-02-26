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

func (m *mockFS) CopyFileRange(cancel <-chan struct{}, in *fuse.CopyFileRangeIn) (written uint32, code fuse.Status) {
	args := m.Called(cancel, in)
	return args.Get(0).(uint32), args.Get(1).(fuse.Status)
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
func (m *mockFS) Link(cancel <-chan struct{}, input *fuse.LinkIn, filename string, out *fuse.EntryOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Symlink(cancel <-chan struct{}, header *fuse.InHeader, pointedTo string, linkName string, out *fuse.EntryOut) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Readlink(cancel <-chan struct{}, header *fuse.InHeader) (out []byte, code fuse.Status) {
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
func (m *mockFS) Release(cancel <-chan struct{}, input *fuse.ReleaseIn)          {}
func (m *mockFS) Write(cancel <-chan struct{}, input *fuse.WriteIn, data []byte) (written uint32, code fuse.Status) {
	return 0, fuse.OK
}
func (m *mockFS) Flush(cancel <-chan struct{}, input *fuse.FlushIn) fuse.Status {
	return fuse.OK
}
func (m *mockFS) Fsync(cancel <-chan struct{}, input *fuse.FsyncIn) fuse.Status {
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
func (m *mockFS) ReleaseDir(input *fuse.ReleaseIn)      {}
func (m *mockFS) FsyncDir(cancel <-chan struct{}, input *fuse.FsyncIn) fuse.Status {
	return fuse.OK
}
func (m *mockFS) StatFs(cancel <-chan struct{}, input *fuse.InHeader, out *fuse.StatfsOut) fuse.Status {
	return fuse.OK
}

func TestCopyFileRange(t *testing.T) {
	mockFs := &mockFS{}
	server := fuse2grpc.NewServer(mockFs)

	tests := []struct {
		name       string
		req        *pb.CopyFileRangeRequest
		mockReturn []interface{}
		want       *pb.CopyFileRangeResponse
		wantErr    error
	}{
		{
			name: "successful copy",
			req: &pb.CopyFileRangeRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
						Pid: 1234,
					},
				},
				FhIn:      2,
				OffIn:     100,
				NodeIdOut: 3,
				FhOut:     4,
				OffOut:    200,
				Len:       1000,
				Flags:     0,
			},
			mockReturn: []interface{}{uint32(1000), fuse.OK},
			want: &pb.CopyFileRangeResponse{
				Written: 1000,
				Status:  &pb.Status{Code: 0},
			},
			wantErr: nil,
		},
		{
			name: "not implemented",
			req: &pb.CopyFileRangeRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
			},
			mockReturn: []interface{}{uint32(0), fuse.ENOSYS},
			want:       nil,
			wantErr:    status.Error(codes.Unimplemented, "method CopyFileRange not implemented"),
		},
		{
			name: "error case",
			req: &pb.CopyFileRangeRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
			},
			mockReturn: []interface{}{uint32(0), fuse.ENOENT},
			want: &pb.CopyFileRangeResponse{
				Written: 0,
				Status:  &pb.Status{Code: int32(fuse.ENOENT)},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFs.On("CopyFileRange", mock.Anything, mock.MatchedBy(func(in *fuse.CopyFileRangeIn) bool {
				return in.NodeId == tt.req.Header.NodeId
			})).Return(tt.mockReturn...).Once()

			got, err := server.CopyFileRange(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockFs.AssertExpectations(t)
		})
	}
}
