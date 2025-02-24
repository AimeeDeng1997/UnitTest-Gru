package fuse2grpc_test

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/chiyutianyi/grpcfuse/fuse2grpc"
	"github.com/chiyutianyi/grpcfuse/pb"
)

// Minimal filesystem implementation
type testFS struct {
	fuse.RawFileSystem
}

func (fs *testFS) Write(cancel <-chan struct{}, in *fuse.WriteIn, data []byte) (written uint32, code fuse.Status) {
	if in.NodeId == 1 {
		return 0, fuse.EIO
	}
	return 50, fuse.OK
}

func TestWrite(t *testing.T) {
	server := fuse2grpc.NewServer(&testFS{})

	tests := []struct {
		name     string
		req      *pb.WriteRequest
		wantResp *pb.WriteResponse
		wantErr  error
	}{
		{
			name: "successful write",
			req: &pb.WriteRequest{
				Header: &pb.InHeader{
					NodeId: 2,
					Caller: &pb.Caller{
						Owner: &pb.Owner{Uid: 1000, Gid: 1000},
						Pid:  12345,
					},
				},
				Fh:         2,
				Offset:     100,
				Size:       50,
				WriteFlags: 0,
				LockOwner:  0,
				Data:       []byte("test data"),
			},
			wantResp: &pb.WriteResponse{
				Written: 50,
				Status:  &pb.Status{Code: 0},
			},
			wantErr: nil,
		},
		{
			name: "write error",
			req: &pb.WriteRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
				Data: []byte("test"),
			},
			wantResp: &pb.WriteResponse{
				Written: 0,
				Status:  &pb.Status{Code: int32(fuse.EIO)},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.Write(context.Background(), tt.req)

			if tt.wantErr != nil && err == nil {
				t.Errorf("Write() error = nil, wantErr %v", tt.wantErr)
				return
			}
			if tt.wantErr != nil && tt.wantErr.Error() != err.Error() {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResp != nil {
				if resp.Written != tt.wantResp.Written {
					t.Errorf("Write() written = %v, want %v", resp.Written, tt.wantResp.Written)
				}
				if resp.Status.Code != tt.wantResp.Status.Code {
					t.Errorf("Write() status = %v, want %v", resp.Status.Code, tt.wantResp.Status.Code)
				}
			}
		})
	}
}
