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
	setXAttrFunc func(cancel <-chan struct{}, in *fuse.SetXAttrIn, attr string, data []byte) fuse.Status
}

func (m *mockFS) SetXAttr(cancel <-chan struct{}, in *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
	if m.setXAttrFunc != nil {
		return m.setXAttrFunc(cancel, in, attr, data)
	}
	return fuse.ENOSYS
}

func TestServer_SetXAttr(t *testing.T) {
	tests := []struct {
		name    string
		fs      fuse.RawFileSystem
		req     *pb.SetXAttrRequest
		want    *pb.SetXAttrResponse
		wantErr error
	}{
		{
			name: "success",
			fs: &mockFS{
				setXAttrFunc: func(cancel <-chan struct{}, in *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
					assert.Equal(t, uint64(123), in.NodeId)
					assert.Equal(t, "test-attr", attr)
					assert.Equal(t, []byte("test-data"), data)
					assert.Equal(t, uint32(9), in.Size)
					assert.Equal(t, uint32(1), in.Flags)
					return fuse.OK
				},
			},
			req: &pb.SetXAttrRequest{
				Header: &pb.InHeader{
					NodeId: 123,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
						Pid: 12345,
					},
				},
				Attr:  "test-attr",
				Data:  []byte("test-data"),
				Size:  9,
				Flags: 1,
			},
			want: &pb.SetXAttrResponse{
				Status: &pb.Status{Code: 0},
			},
		},
		{
			name: "not implemented",
			fs:   &mockFS{},
			req: &pb.SetXAttrRequest{
				Header: &pb.InHeader{
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
				Attr: "test-attr",
			},
			wantErr: status.Errorf(codes.Unimplemented, "method GetXAttr not implemented"),
		},
		{
			name: "permission denied",
			fs: &mockFS{
				setXAttrFunc: func(cancel <-chan struct{}, in *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
					return fuse.EPERM
				},
			},
			req: &pb.SetXAttrRequest{
				Header: &pb.InHeader{
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
				Attr: "test-attr",
			},
			want: &pb.SetXAttrResponse{
				Status: &pb.Status{Code: int32(fuse.EPERM)},
			},
		},
		{
			name: "empty attribute name",
			fs: &mockFS{
				setXAttrFunc: func(cancel <-chan struct{}, in *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
					return fuse.EINVAL
				},
			},
			req: &pb.SetXAttrRequest{
				Header: &pb.InHeader{
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
				Attr: "",
			},
			want: &pb.SetXAttrResponse{
				Status: &pb.Status{Code: int32(fuse.EINVAL)},
			},
		},
		{
			name: "nil data",
			fs: &mockFS{
				setXAttrFunc: func(cancel <-chan struct{}, in *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
					return fuse.OK
				},
			},
			req: &pb.SetXAttrRequest{
				Header: &pb.InHeader{
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
				Attr: "test-attr",
			},
			want: &pb.SetXAttrResponse{
				Status: &pb.Status{Code: 0},
			},
		},
		{
			name: "zero size",
			fs: &mockFS{
				setXAttrFunc: func(cancel <-chan struct{}, in *fuse.SetXAttrIn, attr string, data []byte) fuse.Status {
					assert.Equal(t, uint32(0), in.Size)
					return fuse.OK
				},
			},
			req: &pb.SetXAttrRequest{
				Header: &pb.InHeader{
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
				Attr: "test-attr",
				Data: []byte{},
				Size: 0,
			},
			want: &pb.SetXAttrResponse{
				Status: &pb.Status{Code: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				fs: tt.fs,
			}

			got, err := s.SetXAttr(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
