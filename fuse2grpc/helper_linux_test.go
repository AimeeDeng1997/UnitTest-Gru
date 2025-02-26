package fuse2grpc

import (
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func TestToPbAttr(t *testing.T) {
	tests := []struct {
		name string
		in   *fuse.Attr
		want *pb.Attr
	}{
		{
			name: "normal case",
			in: &fuse.Attr{
				Ino:       123,
				Size:      456,
				Blocks:    789,
				Atime:     1000,
				Mtime:     2000,
				Ctime:     3000,
				Atimensec: 100,
				Mtimensec: 200,
				Ctimensec: 300,
				Mode:      0644,
				Nlink:     1,
				Owner: fuse.Owner{
					Uid: 1000,
					Gid: 1000,
				},
				Rdev:    0,
				Blksize: 4096,
				Padding: 0,
			},
			want: &pb.Attr{
				Ino:       123,
				Size:      456,
				Blocks:    789,
				Atime:     1000,
				Mtime:     2000,
				Ctime:     3000,
				Atimensec: 100,
				Mtimensec: 200,
				Ctimensec: 300,
				Mode:      0644,
				Nlink:     1,
				Owner: &pb.Owner{
					Uid: 1000,
					Gid: 1000,
				},
				Rdev:    0,
				Blksize: 4096,
				Padding: 0,
			},
		},
		{
			name: "zero values",
			in: &fuse.Attr{},
			want: &pb.Attr{
				Owner: &pb.Owner{},
			},
		},
		{
			name: "special values",
			in: &fuse.Attr{
				Mode:  0777,
				Rdev:  1234,
				Nlink: 5,
			},
			want: &pb.Attr{
				Mode:  0777,
				Rdev:  1234,
				Nlink: 5,
				Owner: &pb.Owner{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toPbAttr(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
