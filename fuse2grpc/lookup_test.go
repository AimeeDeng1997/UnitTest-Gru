package fuse2grpc

import (
	"context"
	"testing"

	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockFS struct {
	fuse.RawFileSystem
	lookupFn func(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status
}

func (m *mockFS) Lookup(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
	if m.lookupFn != nil {
		return m.lookupFn(cancel, header, name, out)
	}
	return fuse.ENOSYS
}

func TestLookup(t *testing.T) {
	tests := []struct {
		name     string
		fs       fuse.RawFileSystem
		req      *pb.LookupRequest
		wantResp *pb.LookupResponse
		wantErr  error
	}{
		{
			name: "successful lookup",
			fs: &mockFS{
				lookupFn: func(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
					out.NodeId = 123
					out.Generation = 1
					out.EntryValid = 100
					out.AttrValid = 200
					out.EntryValidNsec = 300
					out.AttrValidNsec = 400
					out.Attr = fuse.Attr{
						Ino:  123,
						Mode: 0644,
						Owner: fuse.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					}
					return fuse.OK
				},
			},
			req: &pb.LookupRequest{
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
			},
			wantResp: &pb.LookupResponse{
				Status: &pb.Status{Code: 0},
				EntryOut: &pb.EntryOut{
					NodeId:         123,
					Generation:     1,
					EntryValid:     100,
					AttrValid:      200,
					EntryValidNsec: 300,
					AttrValidNsec:  400,
					Attr: &pb.Attr{
						Ino:  123,
						Mode: 0644,
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
			},
		},
		{
			name: "not implemented",
			fs:   &mockFS{},
			req: &pb.LookupRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
				Name: "test.txt",
			},
			wantErr: status.Error(codes.Unimplemented, "method Lookup not implemented"),
		},
		{
			name: "lookup error",
			fs: &mockFS{
				lookupFn: func(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
					return fuse.ENOENT
				},
			},
			req: &pb.LookupRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
				Name: "test.txt",
			},
			wantResp: &pb.LookupResponse{
				Status: &pb.Status{Code: int32(fuse.ENOENT)},
			},
		},
		{
			name: "empty name",
			fs: &mockFS{
				lookupFn: func(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
					return fuse.ENOSYS
				},
			},
			req: &pb.LookupRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
				Name: "",
			},
			wantErr: status.Error(codes.Unimplemented, "method Lookup not implemented"),
		},
		{
			name: "large attributes",
			fs: &mockFS{
				lookupFn: func(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
					out.NodeId = 999999
					out.Generation = 88888
					out.EntryValid = 77777
					out.AttrValid = 66666
					out.EntryValidNsec = 55555
					out.AttrValidNsec = 44444
					out.Attr = fuse.Attr{
						Ino:     999999,
						Size:    1<<32 - 1,
						Blocks:  1<<32 - 1,
						Atime:   1<<32 - 1,
						Mtime:   1<<32 - 1,
						Ctime:   1<<32 - 1,
						Mode:    0777,
						Nlink:   100,
						Owner: fuse.Owner{
							Uid: 9999,
							Gid: 9999,
						},
						Rdev:    1234,
						Blksize: 4096,
					}
					return fuse.OK
				},
			},
			req: &pb.LookupRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Name: "big.file",
			},
			wantResp: &pb.LookupResponse{
				Status: &pb.Status{Code: 0},
				EntryOut: &pb.EntryOut{
					NodeId:         999999,
					Generation:     88888,
					EntryValid:     77777,
					AttrValid:      66666,
					EntryValidNsec: 55555,
					AttrValidNsec:  44444,
					Attr: &pb.Attr{
						Ino:       999999,
						Size:      1<<32 - 1,
						Blocks:    1<<32 - 1,
						Atime:     1<<32 - 1,
						Mtime:     1<<32 - 1,
						Ctime:     1<<32 - 1,
						Mode:      0777,
						Nlink:     100,
						Owner: &pb.Owner{
							Uid: 9999,
							Gid: 9999,
						},
						Rdev:    1234,
						Blksize: 4096,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer(tt.fs)
			resp, err := s.Lookup(context.Background(), tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantResp, resp)
		})
	}
}
