package fuse2grpc

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"

	"github.com/chiyutianyi/grpcfuse/pb"
)

type mockRawFS struct {
	fuse.RawFileSystem
	writeFn func(cancel <-chan struct{}, in *fuse.WriteIn, data []byte) (uint32, fuse.Status)
}

func (m *mockRawFS) Write(cancel <-chan struct{}, in *fuse.WriteIn, data []byte) (uint32, fuse.Status) {
	if m.writeFn != nil {
		return m.writeFn(cancel, in, data)
	}
	return 0, fuse.ENOSYS
}

func TestServer_Write(t *testing.T) {
	tests := []struct {
		name    string
		fs      *mockRawFS
		req     *pb.WriteRequest
		want    *pb.WriteResponse
		wantErr bool
	}{
		{
			name: "successful write",
			fs: &mockRawFS{
				writeFn: func(cancel <-chan struct{}, in *fuse.WriteIn, data []byte) (uint32, fuse.Status) {
					return 100, fuse.OK
				},
			},
			req: &pb.WriteRequest{
				Header: &pb.InHeader{
					Length:  100,
					NodeId:  1,
					Caller:  &pb.Caller{Owner: &pb.Owner{Uid: 1000, Gid: 1000}},
					Padding: 0,
				},
				Fh:         123,
				Offset:     0,
				Size:       100,
				WriteFlags: 0,
				LockOwner:  456,
				Data:       []byte("test data"),
			},
			want: &pb.WriteResponse{
				Written: 100,
				Status:  &pb.Status{Code: 0},
			},
			wantErr: false,
		},
		{
			name: "write not implemented",
			fs:   &mockRawFS{},
			req: &pb.WriteRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{Owner: &pb.Owner{Uid: 1000, Gid: 1000}},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "write error",
			fs: &mockRawFS{
				writeFn: func(cancel <-chan struct{}, in *fuse.WriteIn, data []byte) (uint32, fuse.Status) {
					return 0, fuse.EACCES
				},
			},
			req: &pb.WriteRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{Owner: &pb.Owner{Uid: 1000, Gid: 1000}},
				},
				Data: []byte("test data"),
			},
			want: &pb.WriteResponse{
				Written: 0,
				Status:  &pb.Status{Code: int32(fuse.EACCES)},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				fs: tt.fs,
			}

			got, err := s.Write(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
