package grpc2fuse

import (
	"context"
	"testing"

	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mock RawFileSystemClient
type mockRawFileSystemClient struct {
	mock.Mock
	pb.RawFileSystemClient
}

func (m *mockRawFileSystemClient) String(ctx context.Context, in *pb.StringRequest, opts ...grpc.CallOption) (*pb.StringResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.StringResponse), args.Error(1)
}

func TestNewFileSystem(t *testing.T) {
	mockClient := &mockRawFileSystemClient{}
	opts := []grpc.CallOption{grpc.WaitForReady(true)}

	fs := NewFileSystem(mockClient, opts...)

	assert.NotNil(t, fs)
	assert.Equal(t, mockClient, fs.client)
	assert.Equal(t, opts, fs.opts)
}

func TestFileSystem_String(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*mockRawFileSystemClient)
		expected string
	}{
		{
			name: "success",
			setup: func(m *mockRawFileSystemClient) {
				m.On("String", mock.Anything, &pb.StringRequest{}, mock.Anything).
					Return(&pb.StringResponse{Value: "test-fs"}, nil)
			},
			expected: "test-fs",
		},
		{
			name: "grpc error",
			setup: func(m *mockRawFileSystemClient) {
				m.On("String", mock.Anything, &pb.StringRequest{}, mock.Anything).
					Return(&pb.StringResponse{}, status.Error(codes.Internal, "internal error"))
			},
			expected: defaultName,
		},
		{
			name: "empty response",
			setup: func(m *mockRawFileSystemClient) {
				m.On("String", mock.Anything, &pb.StringRequest{}, mock.Anything).
					Return(&pb.StringResponse{}, nil)
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockRawFileSystemClient{}
			tt.setup(mockClient)

			fs := NewFileSystem(mockClient)
			result := fs.String()

			assert.Equal(t, tt.expected, result)
			mockClient.AssertExpectations(t)
		})
	}
}
