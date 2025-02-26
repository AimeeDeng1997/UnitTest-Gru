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
	pb.RawFileSystemClient
}

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
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.MkdirResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Unlink(ctx context.Context, in *pb.UnlinkRequest, opts ...grpc.CallOption) (*pb.UnlinkResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.UnlinkResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Rmdir(ctx context.Context, in *pb.RmdirRequest, opts ...grpc.CallOption) (*pb.RmdirResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.RmdirResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) Rename(ctx context.Context, in *pb.RenameRequest, opts ...grpc.CallOption) (*pb.RenameResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.RenameResponse), args.Error(1)
}

func TestFileSystem_Mkdir(t *testing.T) {
	mockClient := new(MockRawFileSystemClient)
	fs := &fileSystem{
		client: mockClient,
	}

	tests := []struct {
		name     string
		input    *fuse.MkdirIn
		dirname  string
		mockResp *pb.MkdirResponse
		mockErr  error
		want     fuse.Status
	}{
		{
			name: "successful mkdir",
			input: &fuse.MkdirIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Mode:  0755,
				Umask: 0022,
			},
			dirname: "testdir",
			mockResp: &pb.MkdirResponse{
				Status: &pb.Status{Code: 0},
				EntryOut: &pb.EntryOut{
					NodeId: 2,
					Attr: &pb.Attr{
						Ino:   2,
						Mode:  0755,
						Owner: &pb.Owner{},
					},
				},
			},
			mockErr: nil,
			want:    fuse.OK,
		},
		{
			name: "mkdir error",
			input: &fuse.MkdirIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Mode:  0755,
				Umask: 0022,
			},
			dirname:  "testdir",
			mockResp: nil,
			mockErr:  assert.AnError,
			want:     fuse.EIO,
		},
		{
			name: "mkdir permission denied",
			input: &fuse.MkdirIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Mode:  0755,
				Umask: 0022,
			},
			dirname: "testdir",
			mockResp: &pb.MkdirResponse{
				Status: &pb.Status{Code: int32(fuse.EACCES)},
			},
			mockErr: nil,
			want:    fuse.EACCES,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("Mkdir", mock.Anything, &pb.MkdirRequest{
				Header: toPbHeader(&tt.input.InHeader),
				Name:   tt.dirname,
				Mode:   tt.input.Mode,
				Umask:  tt.input.Umask,
			}, mock.Anything).Return(tt.mockResp, tt.mockErr).Once()

			out := &fuse.EntryOut{}
			got := fs.Mkdir(make(chan struct{}), tt.input, tt.dirname, out)
			assert.Equal(t, tt.want, got)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestFileSystem_Unlink(t *testing.T) {
	mockClient := new(MockRawFileSystemClient)
	fs := &fileSystem{
		client: mockClient,
	}

	tests := []struct {
		name     string
		header   *fuse.InHeader
		filename string
		mockResp *pb.UnlinkResponse
		mockErr  error
		want     fuse.Status
	}{
		{
			name: "successful unlink",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			filename: "testfile",
			mockResp: &pb.UnlinkResponse{
				Status: &pb.Status{Code: 0},
			},
			mockErr: nil,
			want:    fuse.OK,
		},
		{
			name: "unlink error",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			filename: "testfile",
			mockResp: nil,
			mockErr:  assert.AnError,
			want:     fuse.EIO,
		},
		{
			name: "unlink not found",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			filename: "testfile",
			mockResp: &pb.UnlinkResponse{
				Status: &pb.Status{Code: int32(fuse.ENOENT)},
			},
			mockErr: nil,
			want:    fuse.ENOENT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("Unlink", mock.Anything, &pb.UnlinkRequest{
				Header: toPbHeader(tt.header),
				Name:   tt.filename,
			}, mock.Anything).Return(tt.mockResp, tt.mockErr).Once()

			got := fs.Unlink(make(chan struct{}), tt.header, tt.filename)
			assert.Equal(t, tt.want, got)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestFileSystem_Rmdir(t *testing.T) {
	mockClient := new(MockRawFileSystemClient)
	fs := &fileSystem{
		client: mockClient,
	}

	tests := []struct {
		name     string
		header   *fuse.InHeader
		dirname  string
		mockResp *pb.RmdirResponse
		mockErr  error
		want     fuse.Status
	}{
		{
			name: "successful rmdir",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			dirname: "testdir",
			mockResp: &pb.RmdirResponse{
				Status: &pb.Status{Code: 0},
			},
			mockErr: nil,
			want:    fuse.OK,
		},
		{
			name: "rmdir error",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			dirname:  "testdir",
			mockResp: nil,
			mockErr:  assert.AnError,
			want:     fuse.EIO,
		},
		{
			name: "rmdir not empty",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			dirname: "testdir",
			mockResp: &pb.RmdirResponse{
				Status: &pb.Status{Code: int32(fuse.ENOENT)},
			},
			mockErr: nil,
			want:    fuse.ENOENT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("Rmdir", mock.Anything, &pb.RmdirRequest{
				Header: toPbHeader(tt.header),
				Name:   tt.dirname,
			}, mock.Anything).Return(tt.mockResp, tt.mockErr).Once()

			got := fs.Rmdir(make(chan struct{}), tt.header, tt.dirname)
			assert.Equal(t, tt.want, got)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestFileSystem_Rename(t *testing.T) {
	mockClient := new(MockRawFileSystemClient)
	fs := &fileSystem{
		client: mockClient,
	}

	tests := []struct {
		name     string
		input    *fuse.RenameIn
		oldName  string
		newName  string
		mockResp *pb.RenameResponse
		mockErr  error
		want     fuse.Status
	}{
		{
			name: "successful rename",
			input: &fuse.RenameIn{
				InHeader: fuse.InHeader{NodeId: 1},
				Newdir:   2,
			},
			oldName: "oldname",
			newName: "newname",
			mockResp: &pb.RenameResponse{
				Status: &pb.Status{Code: 0},
			},
			mockErr: nil,
			want:    fuse.OK,
		},
		{
			name: "rename error",
			input: &fuse.RenameIn{
				InHeader: fuse.InHeader{NodeId: 1},
				Newdir:   2,
			},
			oldName:  "oldname",
			newName:  "newname",
			mockResp: nil,
			mockErr:  assert.AnError,
			want:     fuse.EIO,
		},
		{
			name: "rename target exists",
			input: &fuse.RenameIn{
				InHeader: fuse.InHeader{NodeId: 1},
				Newdir:   2,
			},
			oldName: "oldname",
			newName: "newname",
			mockResp: &pb.RenameResponse{
				Status: &pb.Status{Code: int32(fuse.ENOENT)},
			},
			mockErr: nil,
			want:    fuse.ENOENT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("Rename", mock.Anything, &pb.RenameRequest{
				Header:  toPbHeader(&tt.input.InHeader),
				OldName: tt.oldName,
				NewName: tt.newName,
				Newdir:  tt.input.Newdir,
				Flags:   tt.input.Flags,
				Padding: tt.input.Padding,
			}, mock.Anything).Return(tt.mockResp, tt.mockErr).Once()

			got := fs.Rename(make(chan struct{}), tt.input, tt.oldName, tt.newName)
			assert.Equal(t, tt.want, got)

			mockClient.AssertExpectations(t)
		})
	}
}
