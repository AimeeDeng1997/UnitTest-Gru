package grpc2fuse

import (
	"context"
	"io"
	"testing"

	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MockRawFileSystemClient struct {
	pb.RawFileSystemClient
	OpenDirFunc     func(context.Context, *pb.OpenDirRequest, ...grpc.CallOption) (*pb.OpenDirResponse, error)
	ReadDirFunc     func(context.Context, *pb.ReadDirRequest, ...grpc.CallOption) (pb.RawFileSystem_ReadDirClient, error)
	ReadDirPlusFunc func(context.Context, *pb.ReadDirRequest, ...grpc.CallOption) (pb.RawFileSystem_ReadDirPlusClient, error)
	ReleaseDirFunc  func(context.Context, *pb.ReleaseRequest, ...grpc.CallOption) (*emptypb.Empty, error)
	FsyncDirFunc    func(context.Context, *pb.FsyncRequest, ...grpc.CallOption) (*pb.FsyncResponse, error)
}

func (m *MockRawFileSystemClient) OpenDir(ctx context.Context, req *pb.OpenDirRequest, opts ...grpc.CallOption) (*pb.OpenDirResponse, error) {
	return m.OpenDirFunc(ctx, req, opts...)
}

func (m *MockRawFileSystemClient) ReadDir(ctx context.Context, req *pb.ReadDirRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadDirClient, error) {
	return m.ReadDirFunc(ctx, req, opts...)
}

func (m *MockRawFileSystemClient) ReadDirPlus(ctx context.Context, req *pb.ReadDirRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadDirPlusClient, error) {
	return m.ReadDirPlusFunc(ctx, req, opts...)
}

func (m *MockRawFileSystemClient) ReleaseDir(ctx context.Context, req *pb.ReleaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return m.ReleaseDirFunc(ctx, req, opts...)
}

func (m *MockRawFileSystemClient) FsyncDir(ctx context.Context, req *pb.FsyncRequest, opts ...grpc.CallOption) (*pb.FsyncResponse, error) {
	return m.FsyncDirFunc(ctx, req, opts...)
}

type MockReadDirClient struct {
	pb.RawFileSystem_ReadDirClient
	recvFunc func() (*pb.ReadDirResponse, error)
}

func (m *MockReadDirClient) Recv() (*pb.ReadDirResponse, error) {
	return m.recvFunc()
}

func (m *MockReadDirClient) Header() (metadata.MD, error) {
	return nil, nil
}

func (m *MockReadDirClient) Trailer() metadata.MD {
	return nil
}

func (m *MockReadDirClient) CloseSend() error {
	return nil
}

func (m *MockReadDirClient) Context() context.Context {
	return context.Background()
}

func TestOpenDir(t *testing.T) {
	tests := []struct {
		name     string
		input    *fuse.OpenIn
		mock     func(context.Context, *pb.OpenDirRequest, ...grpc.CallOption) (*pb.OpenDirResponse, error)
		expected fuse.Status
	}{
		{
			name:  "successful open",
			input: &fuse.OpenIn{},
			mock: func(ctx context.Context, req *pb.OpenDirRequest, opts ...grpc.CallOption) (*pb.OpenDirResponse, error) {
				return &pb.OpenDirResponse{
					Status:  &pb.Status{Code: 0},
					OpenOut: &pb.OpenOut{Fh: 1, OpenFlags: 2},
				}, nil
			},
			expected: fuse.OK,
		},
		{
			name:  "grpc error",
			input: &fuse.OpenIn{},
			mock: func(ctx context.Context, req *pb.OpenDirRequest, opts ...grpc.CallOption) (*pb.OpenDirResponse, error) {
				return nil, status.Error(codes.Internal, "internal error")
			},
			expected: fuse.EIO,
		},
		{
			name:  "non-zero status code",
			input: &fuse.OpenIn{},
			mock: func(ctx context.Context, req *pb.OpenDirRequest, opts ...grpc.CallOption) (*pb.OpenDirResponse, error) {
				return &pb.OpenDirResponse{Status: &pb.Status{Code: int32(fuse.ENOENT)}}, nil
			},
			expected: fuse.ENOENT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockRawFileSystemClient{
				OpenDirFunc: tt.mock,
			}
			fs := &fileSystem{client: mockClient}
			out := &fuse.OpenOut{}
			status := fs.OpenDir(make(<-chan struct{}), tt.input, out)
			assert.Equal(t, tt.expected, status)
		})
	}
}

func TestReadDir(t *testing.T) {
	tests := []struct {
		name     string
		input    *fuse.ReadIn
		mock     func() (*pb.ReadDirResponse, error)
		expected fuse.Status
	}{
		{
			name:  "successful read",
			input: &fuse.ReadIn{},
			mock: func() (*pb.ReadDirResponse, error) {
				return &pb.ReadDirResponse{
					Status: &pb.Status{Code: 0},
					Entries: []*pb.DirEntry{
						{Ino: 1, Name: []byte("test"), Mode: 0755},
					},
				}, io.EOF
			},
			expected: fuse.OK,
		},
		{
			name:  "error during stream read",
			input: &fuse.ReadIn{},
			mock: func() (*pb.ReadDirResponse, error) {
				return nil, status.Error(codes.Internal, "stream error")
			},
			expected: fuse.EIO,
		},
		{
			name:  "non-zero status code in response",
			input: &fuse.ReadIn{},
			mock: func() (*pb.ReadDirResponse, error) {
				return &pb.ReadDirResponse{
					Status: &pb.Status{Code: int32(fuse.EACCES)},
				}, nil
			},
			expected: fuse.EACCES,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReadDirClient := &MockReadDirClient{recvFunc: tt.mock}
			mockClient := &MockRawFileSystemClient{
				ReadDirFunc: func(ctx context.Context, req *pb.ReadDirRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadDirClient, error) {
					return mockReadDirClient, nil
				},
			}
			fs := &fileSystem{client: mockClient}
			out := fuse.NewDirEntryList(make([]byte, 1000), 0)
			status := fs.ReadDir(make(<-chan struct{}), tt.input, out)
			assert.Equal(t, tt.expected, status)
		})
	}
}

func TestReadDirPlus(t *testing.T) {
	tests := []struct {
		name     string
		input    *fuse.ReadIn
		mock     func() (*pb.ReadDirResponse, error)
		expected fuse.Status
	}{
		{
			name:  "successful read",
			input: &fuse.ReadIn{},
			mock: func() (*pb.ReadDirResponse, error) {
				return &pb.ReadDirResponse{
					Status: &pb.Status{Code: 0},
					Entries: []*pb.DirEntry{
						{Ino: 1, Name: []byte("test"), Mode: 0755},
					},
				}, io.EOF
			},
			expected: fuse.OK,
		},
		{
			name:  "error case",
			input: &fuse.ReadIn{},
			mock: func() (*pb.ReadDirResponse, error) {
				return nil, status.Error(codes.Internal, "stream error")
			},
			expected: fuse.EIO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReadDirClient := &MockReadDirClient{recvFunc: tt.mock}
			mockClient := &MockRawFileSystemClient{
				ReadDirPlusFunc: func(ctx context.Context, req *pb.ReadDirRequest, opts ...grpc.CallOption) (pb.RawFileSystem_ReadDirPlusClient, error) {
					return mockReadDirClient, nil
				},
			}
			fs := &fileSystem{client: mockClient}
			out := fuse.NewDirEntryList(make([]byte, 1000), 0)
			status := fs.ReadDirPlus(make(<-chan struct{}), tt.input, out)
			assert.Equal(t, tt.expected, status)
		})
	}
}

func TestReleaseDir(t *testing.T) {
	tests := []struct {
		name  string
		input *fuse.ReleaseIn
		mock  func(context.Context, *pb.ReleaseRequest, ...grpc.CallOption) (*emptypb.Empty, error)
	}{
		{
			name:  "successful release",
			input: &fuse.ReleaseIn{},
			mock: func(ctx context.Context, req *pb.ReleaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
				return &emptypb.Empty{}, nil
			},
		},
		{
			name:  "error case",
			input: &fuse.ReleaseIn{},
			mock: func(ctx context.Context, req *pb.ReleaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
				return nil, status.Error(codes.Internal, "release error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockRawFileSystemClient{
				ReleaseDirFunc: tt.mock,
			}
			fs := &fileSystem{client: mockClient}
			fs.ReleaseDir(tt.input)
		})
	}
}

func TestFsyncDir(t *testing.T) {
	tests := []struct {
		name     string
		input    *fuse.FsyncIn
		mock     func(context.Context, *pb.FsyncRequest, ...grpc.CallOption) (*pb.FsyncResponse, error)
		expected fuse.Status
	}{
		{
			name:  "successful fsync",
			input: &fuse.FsyncIn{},
			mock: func(ctx context.Context, req *pb.FsyncRequest, opts ...grpc.CallOption) (*pb.FsyncResponse, error) {
				return &pb.FsyncResponse{Status: &pb.Status{Code: 0}}, nil
			},
			expected: fuse.OK,
		},
		{
			name:  "grpc error",
			input: &fuse.FsyncIn{},
			mock: func(ctx context.Context, req *pb.FsyncRequest, opts ...grpc.CallOption) (*pb.FsyncResponse, error) {
				return nil, status.Error(codes.Internal, "internal error")
			},
			expected: fuse.EIO,
		},
		{
			name:  "non-zero status code",
			input: &fuse.FsyncIn{},
			mock: func(ctx context.Context, req *pb.FsyncRequest, opts ...grpc.CallOption) (*pb.FsyncResponse, error) {
				return &pb.FsyncResponse{Status: &pb.Status{Code: int32(fuse.EACCES)}}, nil
			},
			expected: fuse.EACCES,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockRawFileSystemClient{
				FsyncDirFunc: tt.mock,
			}
			fs := &fileSystem{client: mockClient}
			status := fs.FsyncDir(make(<-chan struct{}), tt.input)
			assert.Equal(t, tt.expected, status)
		})
	}
}
