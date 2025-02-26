package fuse2grpc

import (
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func TestToPbEntryOut(t *testing.T) {
	tests := []struct {
		name string
		in   *fuse.EntryOut
		want *pb.EntryOut
	}{
		{
			name: "basic mapping",
			in: &fuse.EntryOut{
				NodeId:         123,
				Generation:     456,
				EntryValid:     789,
				AttrValid:      321,
				EntryValidNsec: 654,
				AttrValidNsec:  987,
				Attr: fuse.Attr{
					Ino:  111,
					Mode: 0644,
				},
			},
			want: &pb.EntryOut{
				NodeId:         123,
				Generation:     456,
				EntryValid:     789,
				AttrValid:      321,
				EntryValidNsec: 654,
				AttrValidNsec:  987,
				Attr: &pb.Attr{
					Ino:  111,
					Mode: 0644,
				},
			},
		},
		{
			name: "zero values",
			in:   &fuse.EntryOut{},
			want: &pb.EntryOut{
				Attr: &pb.Attr{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toPbEntryOut(tt.in)
			assert.Equal(t, tt.want.NodeId, got.NodeId)
			assert.Equal(t, tt.want.Generation, got.Generation)
			assert.Equal(t, tt.want.EntryValid, got.EntryValid)
			assert.Equal(t, tt.want.AttrValid, got.AttrValid)
			assert.Equal(t, tt.want.EntryValidNsec, got.EntryValidNsec)
			assert.Equal(t, tt.want.AttrValidNsec, got.AttrValidNsec)
			assert.Equal(t, tt.want.Attr.Ino, got.Attr.Ino)
			assert.Equal(t, tt.want.Attr.Mode, got.Attr.Mode)
		})
	}
}

func TestTypeToMode(t *testing.T) {
	tests := []struct {
		name string
		typ  uint32
		want uint32
	}{
		{
			name: "regular file",
			typ:  8,
			want: 0100000,
		},
		{
			name: "directory",
			typ:  4,
			want: 0040000,
		},
		{
			name: "symlink",
			typ:  10,
			want: 0120000,
		},
		{
			name: "zero type",
			typ:  0,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := typeToMode(tt.typ)
			assert.Equal(t, tt.want, got)
		})
	}
}
