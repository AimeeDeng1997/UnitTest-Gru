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

func (m *MockRawFileSystemClient) Lookup(ctx context.Context, in *pb.LookupRequest, opts ...grpc.CallOption) (*pb.LookupResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.LookupResponse), args.Error(1)
}

func TestLookup(t *testing.T) {
	tests := []struct {
		name           string
		header         *fuse.InHeader
		lookupName     string
		mockResponse   *pb.LookupResponse
		mockError      error
		expectedStatus fuse.Status
	}{
		{
			name: "successful lookup",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			lookupName: "testfile",
			mockResponse: &pb.LookupResponse{
				Status: &pb.Status{Code: 0},
				EntryOut: &pb.EntryOut{
					NodeId:     2,
					Generation: 1,
					Attr: &pb.Attr{
						Ino:  2,
						Mode: 0644,
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
			},
			mockError:      nil,
			expectedStatus: fuse.OK,
		},
		{
			name: "lookup not found",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			lookupName: "nonexistent",
			mockResponse: &pb.LookupResponse{
				Status: &pb.Status{Code: int32(fuse.ENOENT)},
			},
			mockError:      nil,
			expectedStatus: fuse.ENOENT,
		},
		{
			name: "grpc error",
			header: &fuse.InHeader{
				NodeId: 1,
			},
			lookupName:     "error",
			mockResponse:   nil,
			mockError:      context.DeadlineExceeded,
			expectedStatus: fuse.EIO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockRawFileSystemClient)
			fs := &fileSystem{
				client: mockClient,
			}

			cancel := make(chan struct{})
			out := &fuse.EntryOut{}

			expectedRequest := &pb.LookupRequest{
				Header: toPbHeader(tt.header),
				Name:   tt.lookupName,
			}

			mockClient.On("Lookup", mock.Anything, expectedRequest, mock.Anything).Return(tt.mockResponse, tt.mockError)

			status := fs.Lookup(cancel, tt.header, tt.lookupName, out)

			assert.Equal(t, tt.expectedStatus, status)
			if status == fuse.OK {
				assert.Equal(t, tt.mockResponse.EntryOut.NodeId, out.NodeId)
				assert.Equal(t, tt.mockResponse.EntryOut.Generation, out.Generation)
				assert.Equal(t, tt.mockResponse.EntryOut.Attr.Mode, out.Attr.Mode)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
