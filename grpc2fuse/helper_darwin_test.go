package grpc2fuse

import (
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func TestGetUmask(t *testing.T) {
	in := &fuse.MknodIn{}
	assert.Equal(t, uint16(0), getUmask(in))
}

func TestSetFlags(t *testing.T) {
	out := &fuse.Attr{}
	flags := uint32(123)
	setFlags(out, flags)
	assert.Equal(t, flags, out.Flags_)
}

func TestSetBlksize(t *testing.T) {
	out := &fuse.Attr{}
	size := uint32(4096)
	setBlksize(out, size)
	// No-op function, just verify it doesn't panic
}

func TestSetPadding(t *testing.T) {
	out := &fuse.Attr{}
	padding := uint32(0)
	setPadding(out, padding)
	// No-op function, just verify it doesn't panic
}

func TestToPbReadIn(t *testing.T) {
	tests := []struct {
		name string
		in   *fuse.ReadIn
		want *pb.ReadIn
	}{
		{
			name: "basic conversion",
			in: &fuse.ReadIn{
				InHeader: fuse.InHeader{
					Length: 100,
					Opcode: 1,
					Unique: 123,
					NodeId: 456,
					Uid:    1000,
					Gid:    1000,
					Pid:    12345,
				},
				Fh:        789,
				Offset:    1000,
				Size:      4096,
				ReadFlags: 1,
			},
			want: &pb.ReadIn{
				Header: &pb.InHeader{
					Length: 100,
					Opcode: 1,
					Unique: 123,
					NodeId: 456,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
						Pid: 12345,
					},
				},
				Fh:        789,
				Offset:    1000,
				Size:      4096,
				ReadFlags: 1,
			},
		},
		{
			name: "zero values",
			in: &fuse.ReadIn{
				InHeader: fuse.InHeader{},
				Fh:      0,
				Offset:  0,
				Size:    0,
			},
			want: &pb.ReadIn{
				Header: &pb.InHeader{
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
				Fh:     0,
				Offset: 0,
				Size:   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toPbReadIn(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
