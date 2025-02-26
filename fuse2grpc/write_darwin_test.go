package fuse2grpc

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/pb"
)

type mockFS struct {
	fuse.RawFileSystem
	writeFunc func(cancel <-chan struct{}, in *fuse.WriteIn, data []byte) (written uint32, code fuse.Status)
}

func (m *mockFS) Write(cancel <-chan struct{}, in *fuse.WriteIn, data []byte) (written uint32, code fuse.Status) {
	if m.writeFunc != nil {
		return m.writeFunc(cancel, in, data)
	}
	return 0, fuse.ENOSYS
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name        string
		fs          *mockFS
		req         *pb.WriteRequest
		wantWritten uint32
		wantStatus  *pb.Status
		wantErr     error
	}{
		{
			name: "successful write",
			fs: &mockFS{
				writeFunc: func(cancel <-chan struct{}, in *fuse.WriteIn, data []byte) (uint32, fuse.Status) {
					return 100, fuse.OK
				},
			},
			req: &pb.WriteRequest{
				Header: &pb.InHeader{
					Length:  100,
					NodeId:  1,
					Unique:  1,
					Opcode:  1,
					Padding: 0,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
						Pid: 1234,
					},
				},
				Fh:         1,
				Offset:     0,
				Size:       100,
				WriteFlags: 0,
				Data:       make([]byte, 100),
			},
			wantWritten: 100,
			wantStatus:  &pb.Status{Code: 0},
			wantErr:     nil,
		},
		{
			name: "write error",
			fs: &mockFS{
				writeFunc: func(cancel <-chan struct{}, in *fuse.WriteIn, data []byte) (uint32, fuse.Status) {
					return 0, fuse.EIO
				},
			},
			req: &pb.WriteRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
				Size: 100,
				Data: make([]byte, 100),
			},
			wantWritten: 0,
			wantStatus:  &pb.Status{Code: int32(fuse.EIO)},
			wantErr:     nil,
		},
		{
			name: "not implemented",
			fs:   &mockFS{},
			req: &pb.WriteRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
			},
			wantWritten: 0,
			wantStatus:  nil,
			wantErr:     status.Errorf(codes.Unimplemented, "method Write not implemented"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{fs: tt.fs}
			resp, err := s.Write(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantWritten, resp.Written)
			assert.Equal(t, tt.wantStatus.Code, resp.Status.Code)
		})
	}
}
