package grpc2fuse

import (
	"context"
	"testing"

	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRawFileSystemClient struct {
	mock.Mock
}

func (m *MockRawFileSystemClient) Create(ctx context.Context, in *pb.CreateRequest, opts ...interface{}) (*pb.CreateResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.CreateResponse), args.Error(1)
}

func TestCreate(t *testing.T) {
	mockClient := new(MockRawFileSystemClient)
	fs := &fileSystem{
		client: mockClient,
	}

	tests := []struct {
		name     string
		input    *fuse.CreateIn
		fileName string
		mockResp *pb.CreateResponse
		mockErr  error
		want     fuse.Status
	}{
		{
			name: "successful create",
			input: &fuse.CreateIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
				Flags: uint32(0644),
				Mode:  uint32(0644),
			},
			fileName: "test.txt",
			mockResp: &pb.CreateResponse{
				Status: &pb.Status{Code: 0},
				EntryOut: &pb.EntryOut{
					NodeId: 2,
					Attr: &pb.Attr{
						Ino:  2,
						Mode: uint32(0644),
					},
				},
				OpenOut: &pb.OpenOut{
					Fh: 1,
				},
			},
			mockErr: nil,
			want:    fuse.OK,
		},
		{
			name: "create error",
			input: &fuse.CreateIn{
				InHeader: fuse.InHeader{
					NodeId: 1,
				},
			},
			fileName: "test.txt",
			mockResp: &pb.CreateResponse{
				Status: &pb.Status{Code: int32(fuse.EACCES)},
			},
			mockErr: nil,
			want:    fuse.EACCES,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On("Create", mock.Anything, &pb.CreateRequest{
				Header: toPbHeader(&tt.input.InHeader),
				Name:   tt.fileName,
				Flags:  tt.input.Flags,
				Mode:   tt.input.Mode,
			}, mock.Anything).Return(tt.mockResp, tt.mockErr)

			out := &fuse.CreateOut{}
			got := fs.Create(make(chan struct{}), tt.input, tt.fileName, out)

			assert.Equal(t, tt.want, got)
			if tt.mockResp != nil && tt.mockResp.Status.Code == 0 {
				assert.Equal(t, tt.mockResp.EntryOut.NodeId, out.EntryOut.NodeId)
				assert.Equal(t, tt.mockResp.OpenOut.Fh, out.OpenOut.Fh)
			}
		})
	}
}
