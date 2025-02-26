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
	mknodFunc func(cancel <-chan struct{}, in *fuse.MknodIn, name string, out *fuse.EntryOut) fuse.Status
}

func (m *mockFS) Mknod(cancel <-chan struct{}, in *fuse.MknodIn, name string, out *fuse.EntryOut) fuse.Status {
	if m.mknodFunc != nil {
		return m.mknodFunc(cancel, in, name, out)
	}
	return fuse.ENOSYS
}

func TestMknod(t *testing.T) {
	tests := []struct {
		name        string
		fs          *mockFS
		req         *pb.MknodRequest
		wantErr     bool
		wantErrCode codes.Code
		wantStatus  int32
		wantEntry   *pb.EntryOut
	}{
		{
			name: "success",
			fs: &mockFS{
				mknodFunc: func(cancel <-chan struct{}, in *fuse.MknodIn, name string, out *fuse.EntryOut) fuse.Status {
					out.NodeId = 123
					out.Generation = 456
					out.EntryValid = 789
					return fuse.OK
				},
			},
			req: &pb.MknodRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Name: "test.txt",
				Mode: 0644,
				Rdev: 0,
			},
			wantStatus: 0,
			wantEntry: &pb.EntryOut{
				NodeId:     123,
				Generation: 456,
				EntryValid: 789,
			},
		},
		{
			name: "not implemented",
			fs:   &mockFS{},
			req: &pb.MknodRequest{
				Header: &pb.InHeader{},
				Name:   "test.txt",
			},
			wantErr:     true,
			wantErrCode: codes.Unimplemented,
		},
		{
			name: "operation failed",
			fs: &mockFS{
				mknodFunc: func(cancel <-chan struct{}, in *fuse.MknodIn, name string, out *fuse.EntryOut) fuse.Status {
					return fuse.EPERM
				},
			},
			req: &pb.MknodRequest{
				Header: &pb.InHeader{},
				Name:   "test.txt",
			},
			wantStatus: int32(fuse.EPERM),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{fs: tt.fs}
			resp, err := s.Mknod(context.Background(), tt.req)

			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.wantErrCode, st.Code())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.Status.Code)

			if tt.wantEntry != nil {
				assert.Equal(t, tt.wantEntry.NodeId, resp.EntryOut.NodeId)
				assert.Equal(t, tt.wantEntry.Generation, resp.EntryOut.Generation)
				assert.Equal(t, tt.wantEntry.EntryValid, resp.EntryOut.EntryValid)
			}
		})
	}
}
