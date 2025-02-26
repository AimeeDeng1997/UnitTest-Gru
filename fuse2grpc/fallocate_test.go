package fuse2grpc

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/pb"
)

type mockRawFileSystem struct {
	mock.Mock
}

func (m *mockRawFileSystem) Init(server *fuse.Server) {}

func (m *mockRawFileSystem) String() string {
	return "mockRawFileSystem"
}

func (m *mockRawFileSystem) SetDebug(debug bool) {}

func (m *mockRawFileSystem) Lookup(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
	return fuse.OK
}

func (m *mockRawFileSystem) Forget(nodeID uint64, nlookup uint64) {}

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

func (m *mockRawFileSystem) Fallocate(cancel <-chan struct{}, in *fuse.FallocateIn) fuse.Status {
	args := m.Called(cancel, in)
	return args.Get(0).(fuse.Status)
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

func (m *mockRawFileSystem) StatFs(cancel <-chan struct{}, input *fuse.InHeader, out *fuse.StatfsOut) fuse.Status {
	return fuse.OK
}

func TestServer_Fallocate(t *testing.T) {
	tests := []struct {
		name    string
		req     *pb.FallocateRequest
		fsResp  fuse.Status
		want    *pb.FallocateResponse
		wantErr error
	}{
		{
			name: "successful fallocate",
			req: &pb.FallocateRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Fh:      123,
				Offset:  0,
				Length:  4096,
				Mode:    0,
				Padding: 0,
			},
			fsResp: fuse.OK,
			want: &pb.FallocateResponse{
				Status: &pb.Status{
					Code: 0,
				},
			},
			wantErr: nil,
		},
		{
			name: "unimplemented fallocate",
			req: &pb.FallocateRequest{
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
			fsResp:  fuse.ENOSYS,
			want:    nil,
			wantErr: status.Errorf(codes.Unimplemented, "method Fallocate not implemented"),
		},
		{
			name: "failed fallocate with error",
			req: &pb.FallocateRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Fh:      123,
				Offset:  0,
				Length:  4096,
				Mode:    0,
				Padding: 0,
			},
			fsResp: fuse.Status(fuse.EIO),
			want: &pb.FallocateResponse{
				Status: &pb.Status{
					Code: int32(fuse.EIO),
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := &mockRawFileSystem{}
			server := NewServer(mockFS)

			mockFS.On("Fallocate", mock.Anything, mock.MatchedBy(func(in *fuse.FallocateIn) bool {
				return in.NodeId == tt.req.Header.NodeId &&
					in.Fh == tt.req.Fh &&
					in.Offset == tt.req.Offset &&
					in.Length == tt.req.Length &&
					in.Mode == tt.req.Mode &&
					in.Padding == tt.req.Padding
			})).Return(tt.fsResp)

			got, err := server.Fallocate(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			mockFS.AssertExpectations(t)
		})
	}
}
