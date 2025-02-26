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

// Mock RawFileSystemClient
type MockRawFileSystemClient struct {
	mock.Mock
	pb.RawFileSystemClient
}

func (m *MockRawFileSystemClient) GetXAttr(ctx context.Context, in *pb.GetXAttrRequest, opts ...grpc.CallOption) (*pb.GetXAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.GetXAttrResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) ListXAttr(ctx context.Context, in *pb.ListXAttrRequest, opts ...grpc.CallOption) (*pb.ListXAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.ListXAttrResponse), args.Error(1)
}

func (m *MockRawFileSystemClient) RemoveXAttr(ctx context.Context, in *pb.RemoveXAttrRequest, opts ...grpc.CallOption) (*pb.RemoveXAttrResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.RemoveXAttrResponse), args.Error(1)
}

func TestGetXAttr(t *testing.T) {
	mockClient := new(MockRawFileSystemClient)
	fs := &fileSystem{
		client: mockClient,
	}

	tests := []struct {
		name        string
		header      *fuse.InHeader
		attr        string
		dest        []byte
		mockResp    *pb.GetXAttrResponse
		mockErr     error
		wantSize    uint32
		wantStatus  fuse.Status
	}{
		{
			name: "success",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			attr: "user.test",
			dest: make([]byte, 10),
			mockResp: &pb.GetXAttrResponse{
				Size: 5,
				Status: &pb.Status{
					Code: 0,
				},
			},
			mockErr:    nil,
			wantSize:   5,
			wantStatus: fuse.OK,
		},
		{
			name: "error_response",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			attr: "user.test",
			dest: make([]byte, 10),
			mockResp: &pb.GetXAttrResponse{
				Status: &pb.Status{
					Code: int32(fuse.ENOENT),
				},
			},
			mockErr:    nil,
			wantSize:   0,
			wantStatus: fuse.ENOENT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("GetXAttr", mock.Anything, &pb.GetXAttrRequest{
				Header: toPbHeader(tt.header),
				Attr:   tt.attr,
				Dest:   tt.dest,
			}, mock.Anything).Return(tt.mockResp, tt.mockErr).Once()

			size, status := fs.GetXAttr(make(chan struct{}), tt.header, tt.attr, tt.dest)
			assert.Equal(t, tt.wantSize, size)
			assert.Equal(t, tt.wantStatus, status)
		})
	}
}

func TestListXAttr(t *testing.T) {
	mockClient := new(MockRawFileSystemClient)
	fs := &fileSystem{
		client: mockClient,
	}

	tests := []struct {
		name        string
		header      *fuse.InHeader
		dest        []byte
		mockResp    *pb.ListXAttrResponse
		mockErr     error
		wantSize    uint32
		wantStatus  fuse.Status
	}{
		{
			name: "success",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			dest: make([]byte, 10),
			mockResp: &pb.ListXAttrResponse{
				Size: 5,
				Status: &pb.Status{
					Code: 0,
				},
			},
			mockErr:    nil,
			wantSize:   5,
			wantStatus: fuse.OK,
		},
		{
			name: "error_response",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			dest: make([]byte, 10),
			mockResp: &pb.ListXAttrResponse{
				Status: &pb.Status{
					Code: int32(fuse.ENOENT),
				},
			},
			mockErr:    nil,
			wantSize:   0,
			wantStatus: fuse.ENOENT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("ListXAttr", mock.Anything, &pb.ListXAttrRequest{
				Header: toPbHeader(tt.header),
				Dest:   tt.dest,
			}, mock.Anything).Return(tt.mockResp, tt.mockErr).Once()

			size, status := fs.ListXAttr(make(chan struct{}), tt.header, tt.dest)
			assert.Equal(t, tt.wantSize, size)
			assert.Equal(t, tt.wantStatus, status)
		})
	}
}

func TestRemoveXAttr(t *testing.T) {
	mockClient := new(MockRawFileSystemClient)
	fs := &fileSystem{
		client: mockClient,
	}

	tests := []struct {
		name        string
		header      *fuse.InHeader
		attr        string
		mockResp    *pb.RemoveXAttrResponse
		mockErr     error
		wantStatus  fuse.Status
	}{
		{
			name: "success",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			attr: "user.test",
			mockResp: &pb.RemoveXAttrResponse{
				Status: &pb.Status{
					Code: 0,
				},
			},
			mockErr:    nil,
			wantStatus: fuse.OK,
		},
		{
			name: "error_response",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			attr: "user.test",
			mockResp: &pb.RemoveXAttrResponse{
				Status: &pb.Status{
					Code: int32(fuse.ENOENT),
				},
			},
			mockErr:    nil,
			wantStatus: fuse.ENOENT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("RemoveXAttr", mock.Anything, &pb.RemoveXAttrRequest{
				Header: toPbHeader(tt.header),
				Attr:   tt.attr,
			}, mock.Anything).Return(tt.mockResp, tt.mockErr).Once()

			status := fs.RemoveXAttr(make(chan struct{}), tt.header, tt.attr)
			assert.Equal(t, tt.wantStatus, status)
		})
	}
}
