package fuse2grpc_test

import (
	"testing"

	"github.com/chiyutianyi/grpcfuse/fuse2grpc"
	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
)

func TestToPbAttr(t *testing.T) {
	tests := []struct {
		name string
		in   *fuse.Attr
		want *pb.Attr
	}{
		{
			name: "empty attr",
			in:   &fuse.Attr{},
			want: &pb.Attr{
				Owner: &pb.Owner{},
			},
		},
		{
			name: "full attr",
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
				Mode:      0755,
				Nlink:     2,
				Uid:       1001,
				Gid:       1002,
				Rdev:      5,
				Flags_:    15,
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
				Mode:      0755,
				Nlink:     2,
				Owner: &pb.Owner{
					Uid: 1001,
					Gid: 1002,
				},
				Rdev:  5,
				Flags: 15,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fuse2grpc.ToPbAttr(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
