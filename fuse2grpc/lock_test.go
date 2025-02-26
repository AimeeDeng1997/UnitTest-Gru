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

func (m *MockRawFileSystem) Init(*fuse.Server) {
	m.Called()
}

func (m *MockRawFileSystem) String() string {
	return "MockRawFileSystem"
}

func (m *MockRawFileSystem) SetDebug(debug bool) {
	m.Called(debug)
}

func (m *MockRawFileSystem) StatFs(cancel <-chan struct{}, in *fuse.InHeader, out *fuse.StatfsOut) fuse.Status {
	args := m.Called(cancel, in, out)
	return args.Get(0).(fuse.Status)
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

func (m *MockRawFileSystem) Link(cancel <-chan struct{}, input *fuse.LinkIn, name string, out *fuse.EntryOut) fuse.Status {
	args := m.Called(cancel, input, name, out)
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

func (m *MockRawFileSystem) GetXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string, dest []byte) (uint32, fuse.Status) {
	args := m.Called(cancel, header, attr, dest)
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

func (m *MockRawFileSystem) Lseek(cancel <-chan struct{}, in *fuse.LseekIn, out *fuse.LseekOut) fuse.Status {
	args := m.Called(cancel, in, out)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) GetLk(cancel <-chan struct{}, in *fuse.LkIn, out *fuse.LkOut) fuse.Status {
	args := m.Called(cancel, in, out)
	if args.Get(0) == fuse.OK && out != nil && args.Get(1) != nil {
		*out = *(args.Get(1).(*fuse.LkOut))
	}
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) SetLk(cancel <-chan struct{}, in *fuse.LkIn) fuse.Status {
	args := m.Called(cancel, in)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) SetLkw(cancel <-chan struct{}, in *fuse.LkIn) fuse.Status {
	args := m.Called(cancel, in)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Release(cancel <-chan struct{}, input *fuse.ReleaseIn) {
	m.Called(cancel, input)
}

func (m *MockRawFileSystem) Write(cancel <-chan struct{}, input *fuse.WriteIn, data []byte) (uint32, fuse.Status) {
	args := m.Called(cancel, input, data)
	return args.Get(0).(uint32), args.Get(1).(fuse.Status)
}

func (m *MockRawFileSystem) CopyFileRange(cancel <-chan struct{}, input *fuse.CopyFileRangeIn) (uint32, fuse.Status) {
	args := m.Called(cancel, input)
	return args.Get(0).(uint32), args.Get(1).(fuse.Status)
}

func (m *MockRawFileSystem) Flush(cancel <-chan struct{}, input *fuse.FlushIn) fuse.Status {
	args := m.Called(cancel, input)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Fsync(cancel <-chan struct{}, input *fuse.FsyncIn) fuse.Status {
	args := m.Called(cancel, input)
	return args.Get(0).(fuse.Status)
}

func (m *MockRawFileSystem) Fallocate(cancel <-chan struct{}, input *fuse.FallocateIn) fuse.Status {
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

func TestGetLk(t *testing.T) {
	mockFS := &MockRawFileSystem{}
	server := fuse2grpc.NewServer(mockFS)

	tests := []struct {
		name     string
		req      *pb.LkRequest
		mockResp fuse.Status
		mockOut  *fuse.LkOut
		wantErr  error
		wantResp *pb.GetLkResponse
	}{
		{
			name: "successful get lock",
			req: &pb.LkRequest{
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
				Fh:    123,
				Owner: 456,
				Lk: &pb.FileLock{
					Start: 0,
					End:   100,
					Type:  1,
					Pid:   789,
				},
			},
			mockResp: fuse.OK,
			mockOut: &fuse.LkOut{
				Lk: fuse.FileLock{
					Start: 0,
					End:   100,
					Typ:   1,
					Pid:   789,
				},
			},
			wantResp: &pb.GetLkResponse{
				Status: &pb.Status{Code: 0},
				Lk: &pb.FileLock{
					Start: 0,
					End:   100,
					Type:  1,
					Pid:   789,
				},
			},
		},
		{
			name: "not implemented",
			req: &pb.LkRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Lk: &pb.FileLock{},
			},
			mockResp: fuse.ENOSYS,
			mockOut:  nil,
			wantErr:  status.Errorf(codes.Unimplemented, "method GetLk not implemented"),
		},
		{
			name: "error status",
			req: &pb.LkRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Lk: &pb.FileLock{},
			},
			mockResp: fuse.EACCES,
			mockOut:  &fuse.LkOut{},
			wantResp: &pb.GetLkResponse{
				Status: &pb.Status{Code: int32(fuse.EACCES)},
				Lk:    &pb.FileLock{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS.On("GetLk", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResp, tt.mockOut).Once()

			got, err := server.GetLk(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantResp, got)
			mockFS.AssertExpectations(t)
		})
	}
}

func TestSetLk(t *testing.T) {
	mockFS := &MockRawFileSystem{}
	server := fuse2grpc.NewServer(mockFS)

	tests := []struct {
		name     string
		req      *pb.LkRequest
		mockResp fuse.Status
		wantErr  error
		wantResp *pb.SetLkResponse
	}{
		{
			name: "successful set lock",
			req: &pb.LkRequest{
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
				Fh:    123,
				Owner: 456,
				Lk: &pb.FileLock{
					Start: 0,
					End:   100,
					Type:  1,
					Pid:   789,
				},
			},
			mockResp: fuse.OK,
			wantResp: &pb.SetLkResponse{
				Status: &pb.Status{Code: 0},
			},
		},
		{
			name: "not implemented",
			req: &pb.LkRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Lk: &pb.FileLock{},
			},
			mockResp: fuse.ENOSYS,
			wantErr:  status.Errorf(codes.Unimplemented, "method SetLk not implemented"),
		},
		{
			name: "error status",
			req: &pb.LkRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Lk: &pb.FileLock{},
			},
			mockResp: fuse.EACCES,
			wantResp: &pb.SetLkResponse{
				Status: &pb.Status{Code: int32(fuse.EACCES)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS.On("SetLk", mock.Anything, mock.Anything).Return(tt.mockResp).Once()

			got, err := server.SetLk(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantResp, got)
			mockFS.AssertExpectations(t)
		})
	}
}

func TestSetLkw(t *testing.T) {
	mockFS := &MockRawFileSystem{}
	server := fuse2grpc.NewServer(mockFS)

	tests := []struct {
		name     string
		req      *pb.LkRequest
		mockResp fuse.Status
		wantErr  error
		wantResp *pb.SetLkResponse
	}{
		{
			name: "successful set lock wait",
			req: &pb.LkRequest{
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
				Fh:    123,
				Owner: 456,
				Lk: &pb.FileLock{
					Start: 0,
					End:   100,
					Type:  1,
					Pid:   789,
				},
			},
			mockResp: fuse.OK,
			wantResp: &pb.SetLkResponse{
				Status: &pb.Status{Code: 0},
			},
		},
		{
			name: "not implemented",
			req: &pb.LkRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Lk: &pb.FileLock{},
			},
			mockResp: fuse.ENOSYS,
			wantErr:  status.Errorf(codes.Unimplemented, "method SetLkw not implemented"),
		},
		{
			name: "error status",
			req: &pb.LkRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Lk: &pb.FileLock{},
			},
			mockResp: fuse.EACCES,
			wantResp: &pb.SetLkResponse{
				Status: &pb.Status{Code: int32(fuse.EACCES)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS.On("SetLk", mock.Anything, mock.Anything).Return(tt.mockResp).Once()

			got, err := server.SetLkw(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantResp, got)
			mockFS.AssertExpectations(t)
		})
	}
}
