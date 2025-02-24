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
	lookupFunc func(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status
}

func (m *mockFS) Lookup(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
	if m.lookupFunc != nil {
		return m.lookupFunc(cancel, header, name, out)
	}
	return fuse.ENOSYS
}

func TestLookup(t *testing.T) {
	tests := []struct {
		name     string
		fs       fuse.RawFileSystem
		req      *pb.LookupRequest
		expected *pb.LookupResponse
		err      error
	}{
		{
			name: "successful lookup",
			fs: &mockFS{
				lookupFunc: func(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
					out.NodeId = 123
					out.Generation = 1
					out.EntryValid = 1000
					out.AttrValid = 1000
					out.EntryValidNsec = 0
					out.AttrValidNsec = 0
					out.Attr = fuse.Attr{
						Ino:       123,
						Size:      1024,
						Blocks:    2,
						Atime:     1000,
						Mtime:     1000,
						Ctime:     1000,
						Atimensec: 0,
						Mtimensec: 0,
						Ctimensec: 0,
						Mode:      0644,
						Nlink:     1,
						Owner: fuse.Owner{
							Uid: header.Uid,
							Gid: header.Gid,
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
						Pid: 12345,
					},
				},
				Name: "test.txt",
			},
			expected: &pb.LookupResponse{
				Status: &pb.Status{Code: 0},
				EntryOut: &pb.EntryOut{
					NodeId:         123,
					Generation:     1,
					EntryValid:     1000,
					AttrValid:      1000,
					EntryValidNsec: 0,
					AttrValidNsec:  0,
					Attr: &pb.Attr{
						Ino:       123,
						Size:      1024,
						Blocks:    2,
						Atime:     1000,
						Mtime:     1000,
						Ctime:     1000,
						Atimensec: 0,
						Mtimensec: 0,
						Ctimensec: 0,
						Mode:      0644,
						Nlink:     1,
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "not implemented",
			fs:   &mockFS{},
			req: &pb.LookupRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
						Pid: 12345,
					},
				},
				Name: "test.txt",
			},
			expected: nil,
			err:      status.Error(codes.Unimplemented, "method Lookup not implemented"),
		},
		{
			name: "lookup error",
			fs: &mockFS{
				lookupFunc: func(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
					return fuse.ENOENT
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
						Pid: 12345,
					},
				},
				Name: "nonexistent.txt",
			},
			expected: &pb.LookupResponse{
				Status: &pb.Status{Code: int32(fuse.ENOENT)},
			},
			err: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := NewServer(tc.fs)
			resp, err := server.Lookup(context.Background(), tc.req)

			if tc.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, resp)
		})
	}
}
