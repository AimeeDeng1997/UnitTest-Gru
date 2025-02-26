package grpc2fuse

import (
	"context"
	"testing"

	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockRawFileSystemClient struct {
	mock.Mock
	pb.RawFileSystemClient
}

func (m *MockRawFileSystemClient) GetLk(ctx context.Context, in *pb.LkRequest, opts ...grpc.CallOption) (*pb.GetLkResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.GetLkResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) SetLk(ctx context.Context, in *pb.LkRequest, opts ...grpc.CallOption) (*pb.SetLkResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.SetLkResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) SetLkw(ctx context.Context, in *pb.LkRequest, opts ...grpc.CallOption) (*pb.SetLkResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.SetLkResponse), args.Error(1)
}

func TestFileSystem_GetLk(t *testing.T) {
	mockClient := new(MockRawFileSystemClient)
	fs := &fileSystem{
		client: mockClient,
	}

	tests := []struct {
		name     string
		input    *fuse.LkIn
		mockResp *pb.GetLkResponse
		mockErr  error
		want     fuse.Status
	}{
		{
			name: "successful get lock",
			input: &fuse.LkIn{
				InHeader: fuse.InHeader{},
				Fh:      1,
				Owner:   2,
				Lk: fuse.FileLock{
					Start: 0,
					End:   100,
					Typ:   1,
					Pid:   123,
				},
			},
			mockResp: &pb.GetLkResponse{
				Status: &pb.Status{Code: 0},
				Lk: &pb.FileLock{
					Start: 0,
					End:   100,
					Type:  1,
					Pid:   123,
				},
			},
			mockErr: nil,
			want:    fuse.OK,
		},
		{
			name: "error response",
			input: &fuse.LkIn{
				InHeader: fuse.InHeader{},
				Fh:      1,
			},
			mockResp: &pb.GetLkResponse{
				Status: &pb.Status{Code: int32(fuse.EACCES)},
			},
			mockErr: nil,
			want:    fuse.EACCES,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("GetLk", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResp, tt.mockErr).Once()

			out := &fuse.LkOut{}
			got := fs.GetLk(make(chan struct{}), tt.input, out)

			assert.Equal(t, tt.want, got)
			if tt.mockResp != nil && tt.mockResp.Status.Code == 0 && tt.mockResp.Lk != nil {
				assert.Equal(t, tt.mockResp.Lk.Start, out.Lk.Start)
				assert.Equal(t, tt.mockResp.Lk.End, out.Lk.End)
				assert.Equal(t, tt.mockResp.Lk.Type, out.Lk.Typ)
				assert.Equal(t, tt.mockResp.Lk.Pid, out.Lk.Pid)
			}
		})
	}
}

func TestFileSystem_SetLk(t *testing.T) {
	mockClient := new(MockRawFileSystemClient)
	fs := &fileSystem{
		client: mockClient,
	}

	tests := []struct {
		name     string
		input    *fuse.LkIn
		mockResp *pb.SetLkResponse
		mockErr  error
		want     fuse.Status
	}{
		{
			name: "successful set lock",
			input: &fuse.LkIn{
				InHeader: fuse.InHeader{},
				Fh:      1,
				Owner:   2,
				Lk: fuse.FileLock{
					Start: 0,
					End:   100,
					Typ:   1,
					Pid:   123,
				},
			},
			mockResp: &pb.SetLkResponse{
				Status: &pb.Status{Code: 0},
			},
			mockErr: nil,
			want:    fuse.OK,
		},
		{
			name: "error response",
			input: &fuse.LkIn{
				InHeader: fuse.InHeader{},
				Fh:      1,
			},
			mockResp: &pb.SetLkResponse{
				Status: &pb.Status{Code: int32(fuse.EACCES)},
			},
			mockErr: nil,
			want:    fuse.EACCES,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("SetLk", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResp, tt.mockErr).Once()

			got := fs.SetLk(make(chan struct{}), tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFileSystem_SetLkw(t *testing.T) {
	mockClient := new(MockRawFileSystemClient)
	fs := &fileSystem{
		client: mockClient,
	}

	tests := []struct {
		name     string
		input    *fuse.LkIn
		mockResp *pb.SetLkResponse
		mockErr  error
		want     fuse.Status
	}{
		{
			name: "successful set lock wait",
			input: &fuse.LkIn{
				InHeader: fuse.InHeader{},
				Fh:      1,
				Owner:   2,
				Lk: fuse.FileLock{
					Start: 0,
					End:   100,
					Typ:   1,
					Pid:   123,
				},
			},
			mockResp: &pb.SetLkResponse{
				Status: &pb.Status{Code: 0},
			},
			mockErr: nil,
			want:    fuse.OK,
		},
		{
			name: "error response",
			input: &fuse.LkIn{
				InHeader: fuse.InHeader{},
				Fh:      1,
			},
			mockResp: &pb.SetLkResponse{
				Status: &pb.Status{Code: int32(fuse.EACCES)},
			},
			mockErr: nil,
			want:    fuse.EACCES,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("SetLkw", mock.Anything, mock.Anything, mock.Anything).Return(tt.mockResp, tt.mockErr).Once()

			got := fs.SetLkw(make(chan struct{}), tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
