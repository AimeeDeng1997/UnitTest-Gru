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

type MockRawFileSystem struct {
	mock.Mock
}

func (m *MockRawFileSystem) String() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRawFileSystem) SetDebug(debug bool) {}

func (m *MockRawFileSystem) Init(server *fuse.Server) {
	m.Called(server)
}

func (m *MockRawFileSystem) Lookup(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
	args := m.Called(cancel, header, name, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Forget(nodeID uint64, nlookup uint64) {
	m.Called(nodeID, nlookup)
}

func (m *MockRawFileSystem) GetAttr(cancel <-chan struct{}, input *fuse.GetAttrIn, out *fuse.AttrOut) fuse.Status {
	args := m.Called(cancel, input, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) SetAttr(cancel <-chan struct{}, input *fuse.SetAttrIn, out *fuse.AttrOut) fuse.Status {
	args := m.Called(cancel, input, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Mknod(cancel <-chan struct{}, input *fuse.MknodIn, name string, out *fuse.EntryOut) fuse.Status {
	args := m.Called(cancel, input, name, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Mkdir(cancel <-chan struct{}, input *fuse.MkdirIn, name string, out *fuse.EntryOut) fuse.Status {
	args := m.Called(cancel, input, name, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Unlink(cancel <-chan struct{}, header *fuse.InHeader, name string) fuse.Status {
	args := m.Called(cancel, header, name)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Rmdir(cancel <-chan struct{}, header *fuse.InHeader, name string) fuse.Status {
	args := m.Called(cancel, header, name)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Rename(cancel <-chan struct{}, input *fuse.RenameIn, oldName string, newName string) fuse.Status {
	args := m.Called(cancel, input, oldName, newName)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Link(cancel <-chan struct{}, input *fuse.LinkIn, filename string, out *fuse.EntryOut) fuse.Status {
	args := m.Called(cancel, input, filename, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Symlink(cancel <-chan struct{}, header *fuse.InHeader, pointedTo string, linkName string, out *fuse.EntryOut) fuse.Status {
	args := m.Called(cancel, header, pointedTo, linkName, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Readlink(cancel <-chan struct{}, header *fuse.InHeader) ([]byte, fuse.Status) {
	args := m.Called(cancel, header)
	return args.Get(0).([]byte), args.Get(1).(fuse.Status)
}

func (m *MockRawFileSystem) Access(cancel <-chan struct{}, input *fuse.AccessIn) fuse.Status {
	args := m.Called(cancel, input)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) GetXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string, data []byte) (uint32, fuse.Status) {
	args := m.Called(cancel, header, attr, data)
	return args.Get(0).(uint32), args.Get(1).(fuse.Status)
}

func (m *MockRawFileSystem) ListXAttr(cancel <-chan struct{}, header *fuse.InHeader, dest []byte) (uint32, fuse.Status) {
	args := m.Called(cancel, header, dest)
	return args.Get(0).(uint32), args.Get(1).(fuse.Status)
}

func (m *MockRawFileSystem) SetXAttr(cancel <-chan struct{}, input *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
	args := m.Called(cancel, input, attr, data)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) RemoveXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string) fuse.Status {
	args := m.Called(cancel, header, attr)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Create(cancel <-chan struct{}, input *fuse.CreateIn, name string, out *fuse.CreateOut) fuse.Status {
	args := m.Called(cancel, input, name, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Open(cancel <-chan struct{}, input *fuse.OpenIn, out *fuse.OpenOut) fuse.Status {
	args := m.Called(cancel, input, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Read(cancel <-chan struct{}, input *fuse.ReadIn, buf []byte) (fuse.ReadResult, fuse.Status) {
	args := m.Called(cancel, input, buf)
	return args.Get(0).(fuse.ReadResult), args.Get(1).(fuse.Status)
}

func (m *MockRawFileSystem) Write(cancel <-chan struct{}, input *fuse.WriteIn, data []byte) (uint32, fuse.Status) {
	args := m.Called(cancel, input, data)
	return args.Get(0).(uint32), args.Get(1).(fuse.Status)
}

func (m *MockRawFileSystem) Flush(cancel <-chan struct{}, input *fuse.FlushIn) fuse.Status {
	args := m.Called(cancel, input)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Release(cancel <-chan struct{}, input *fuse.ReleaseIn) {
	m.Called(cancel, input)
}

func (m *MockRawFileSystem) Fsync(cancel <-chan struct{}, input *fuse.FsyncIn) fuse.Status {
	args := m.Called(cancel, input)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) OpenDir(cancel <-chan struct{}, input *fuse.OpenIn, out *fuse.OpenOut) fuse.Status {
	args := m.Called(cancel, input, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) ReadDir(cancel <-chan struct{}, input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	args := m.Called(cancel, input, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) ReadDirPlus(cancel <-chan struct{}, input *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	args := m.Called(cancel, input, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) ReleaseDir(input *fuse.ReleaseIn) {
	m.Called(input)
}

func (m *MockRawFileSystem) FsyncDir(cancel <-chan struct{}, input *fuse.FsyncIn) fuse.Status {
	args := m.Called(cancel, input)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) StatFs(cancel <-chan struct{}, input *fuse.InHeader, out *fuse.StatfsOut) fuse.Status {
	args := m.Called(cancel, input, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) CopyFileRange(cancel <-chan struct{}, input *fuse.CopyFileRangeIn) (uint32, fuse.Status) {
	args := m.Called(cancel, input)
	return args.Get(0).(uint32), args.Get(1).(fuse.Status)
}

func (m *MockRawFileSystem) Lseek(cancel <-chan struct{}, in *fuse.LseekIn, out *fuse.LseekOut) fuse.Status {
	args := m.Called(cancel, in, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) GetLk(cancel <-chan struct{}, input *fuse.LkIn, out *fuse.LkOut) fuse.Status {
	args := m.Called(cancel, input, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) SetLk(cancel <-chan struct{}, input *fuse.LkIn) fuse.Status {
	args := m.Called(cancel, input)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) SetLkw(cancel <-chan struct{}, input *fuse.LkIn) fuse.Status {
	args := m.Called(cancel, input)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Fallocate(cancel <-chan struct{}, input *fuse.FallocateIn) fuse.Status {
	args := m.Called(cancel, input)
	return args.Get(0).(fuse.Status)
}

func TestMkdir(t *testing.T) {
	mockFS := new(MockRawFileSystem)
	server := fuse2grpc.NewServer(mockFS)

	tests := []struct {
		name     string
		request  *pb.MkdirRequest
		fsStatus fuse.Status
		wantErr  error
	}{
		{
			name: "successful mkdir",
			request: &pb.MkdirRequest{
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
				Name:  "testdir",
				Mode:  0755,
				Umask: 0022,
			},
			fsStatus: fuse.OK,
			wantErr:  nil,
		},
		{
			name: "unimplemented mkdir",
			request: &pb.MkdirRequest{
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
				Name:  "testdir",
				Mode:  0755,
				Umask: 0022,
			},
			fsStatus: fuse.ENOSYS,
			wantErr:  status.Errorf(codes.Unimplemented, "method Mkdir not implemented"),
		},
		{
			name: "permission denied mkdir",
			request: &pb.MkdirRequest{
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
				Name:  "testdir",
				Mode:  0755,
				Umask: 0022,
			},
			fsStatus: fuse.EPERM,
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS.On("Mkdir", mock.Anything, mock.MatchedBy(func(input *fuse.MkdirIn) bool {
				return input.NodeId == tt.request.Header.NodeId &&
					input.Mode == tt.request.Mode &&
					input.Umask == tt.request.Umask
			}), tt.request.Name, mock.Anything).Return(tt.fsStatus).Once()

			resp, err := server.Mkdir(context.Background(), tt.request)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, int32(tt.fsStatus), resp.Status.Code)
			mockFS.AssertExpectations(t)
		})
	}
}

func TestUnlink(t *testing.T) {
	mockFS := new(MockRawFileSystem)
	server := fuse2grpc.NewServer(mockFS)

	tests := []struct {
		name     string
		request  *pb.UnlinkRequest
		fsStatus fuse.Status
		wantErr  error
	}{
		{
			name: "successful unlink",
			request: &pb.UnlinkRequest{
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
				Name: "testfile",
			},
			fsStatus: fuse.OK,
			wantErr:  nil,
		},
		{
			name: "unimplemented unlink",
			request: &pb.UnlinkRequest{
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
				Name: "testfile",
			},
			fsStatus: fuse.ENOSYS,
			wantErr:  status.Errorf(codes.Unimplemented, "method Unlink not implemented"),
		},
		{
			name: "no such file unlink",
			request: &pb.UnlinkRequest{
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
				Name: "nonexistent",
			},
			fsStatus: fuse.ENOENT,
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS.On("Unlink", mock.Anything, mock.MatchedBy(func(header *fuse.InHeader) bool {
				return header.NodeId == tt.request.Header.NodeId
			}), tt.request.Name).Return(tt.fsStatus).Once()

			resp, err := server.Unlink(context.Background(), tt.request)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, int32(tt.fsStatus), resp.Status.Code)
			mockFS.AssertExpectations(t)
		})
	}
}

func TestRmdir(t *testing.T) {
	mockFS := new(MockRawFileSystem)
	server := fuse2grpc.NewServer(mockFS)

	tests := []struct {
		name     string
		request  *pb.RmdirRequest
		fsStatus fuse.Status
		wantErr  error
	}{
		{
			name: "successful rmdir",
			request: &pb.RmdirRequest{
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
				Name: "testdir",
			},
			fsStatus: fuse.OK,
			wantErr:  nil,
		},
		{
			name: "unimplemented rmdir",
			request: &pb.RmdirRequest{
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
				Name: "testdir",
			},
			fsStatus: fuse.ENOSYS,
			wantErr:  status.Errorf(codes.Unimplemented, "method Rmdir not implemented"),
		},
		{
			name: "directory not empty rmdir",
			request: &pb.RmdirRequest{
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
				Name: "nonempty",
			},
			fsStatus: fuse.EBUSY,
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS.On("Rmdir", mock.Anything, mock.MatchedBy(func(header *fuse.InHeader) bool {
				return header.NodeId == tt.request.Header.NodeId
			}), tt.request.Name).Return(tt.fsStatus).Once()

			resp, err := server.Rmdir(context.Background(), tt.request)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, int32(tt.fsStatus), resp.Status.Code)
			mockFS.AssertExpectations(t)
		})
	}
}

func TestRename(t *testing.T) {
	mockFS := new(MockRawFileSystem)
	server := fuse2grpc.NewServer(mockFS)

	tests := []struct {
		name     string
		request  *pb.RenameRequest
		fsStatus fuse.Status
		wantErr  error
	}{
		{
			name: "successful rename",
			request: &pb.RenameRequest{
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
				OldName: "oldname",
				NewName: "newname",
				Newdir:  2,
				Flags:   0,
			},
			fsStatus: fuse.OK,
			wantErr:  nil,
		},
		{
			name: "unimplemented rename",
			request: &pb.RenameRequest{
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
				OldName: "oldname",
				NewName: "newname",
				Newdir:  2,
				Flags:   0,
			},
			fsStatus: fuse.ENOSYS,
			wantErr:  status.Errorf(codes.Unimplemented, "method Rename not implemented"),
		},
		{
			name: "target exists rename",
			request: &pb.RenameRequest{
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
				OldName: "oldname",
				NewName: "existing",
				Newdir:  2,
				Flags:   0,
			},
			fsStatus: fuse.EBUSY,
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS.On("Rename", mock.Anything, mock.MatchedBy(func(input *fuse.RenameIn) bool {
				return input.NodeId == tt.request.Header.NodeId &&
					input.Newdir == tt.request.Newdir &&
					input.Flags == tt.request.Flags
			}), tt.request.OldName, tt.request.NewName).Return(tt.fsStatus).Once()

			resp, err := server.Rename(context.Background(), tt.request)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, int32(tt.fsStatus), resp.Status.Code)
			mockFS.AssertExpectations(t)
		})
	}
}
