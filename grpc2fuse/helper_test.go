package grpc2fuse

import (
	"testing"

	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestToPbHeader(t *testing.T) {
	header := &fuse.InHeader{
		Length: 100,
		Opcode: 1,
		Unique: 123,
		NodeId: 456,
	}

	result := toPbHeader(header)

	assert.Equal(t, uint32(100), result.Length)
	assert.Equal(t, uint32(1), result.Opcode)
	assert.Equal(t, uint64(123), result.Unique)
	assert.Equal(t, uint64(456), result.NodeId)
	assert.NotNil(t, result.Caller)
	assert.NotNil(t, result.Caller.Owner)
}

func TestToFuseAttr(t *testing.T) {
	in := &pb.Attr{
		Ino:       1,
		Size:      1000,
		Blocks:    2,
		Atime:     123456,
		Mtime:     234567,
		Ctime:     345678,
		Atimensec: 1,
		Mtimensec: 2,
		Ctimensec: 3,
		Mode:      0644,
		Nlink:     1,
		Owner: &pb.Owner{
			Uid: 1000,
			Gid: 1000,
		},
		Rdev:    0,
		Flags:   0,
		Blksize: 4096,
		Padding: 0,
	}

	out := &fuse.Attr{}
	toFuseAttr(out, in)

	assert.Equal(t, uint64(1), out.Ino)
	assert.Equal(t, uint64(1000), out.Size)
	assert.Equal(t, uint64(2), out.Blocks)
	assert.Equal(t, uint64(123456), out.Atime)
	assert.Equal(t, uint64(234567), out.Mtime)
	assert.Equal(t, uint64(345678), out.Ctime)
	assert.Equal(t, uint32(1), out.Atimensec)
	assert.Equal(t, uint32(2), out.Mtimensec)
	assert.Equal(t, uint32(3), out.Ctimensec)
	assert.Equal(t, uint32(0644), out.Mode)
	assert.Equal(t, uint32(1), out.Nlink)
	assert.Equal(t, uint32(1000), out.Uid)
	assert.Equal(t, uint32(1000), out.Gid)
	assert.Equal(t, uint32(0), out.Rdev)
}

func TestToFuseAttrNilOwner(t *testing.T) {
	in := &pb.Attr{
		Ino:    1,
		Size:   1000,
		Blocks: 2,
		Owner: &pb.Owner{}, // Empty owner instead of nil
	}

	out := &fuse.Attr{}
	toFuseAttr(out, in)

	assert.Equal(t, uint64(1), out.Ino)
	assert.Equal(t, uint64(1000), out.Size)
	assert.Equal(t, uint64(2), out.Blocks)
	assert.Equal(t, uint32(0), out.Uid)
	assert.Equal(t, uint32(0), out.Gid)
}

func TestToFuseEntryOut(t *testing.T) {
	in := &pb.EntryOut{
		NodeId:         123,
		Generation:     1,
		AttrValid:      3600,
		AttrValidNsec:  0,
		EntryValid:     3600,
		EntryValidNsec: 0,
		Attr: &pb.Attr{
			Ino:  1,
			Size: 1000,
			Owner: &pb.Owner{
				Uid: 1000,
				Gid: 1000,
			},
		},
	}

	out := &fuse.EntryOut{}
	toFuseEntryOut(out, in)

	assert.Equal(t, uint64(123), out.NodeId)
	assert.Equal(t, uint64(1), out.Generation)
	assert.Equal(t, uint64(3600), out.AttrValid)
	assert.Equal(t, uint32(0), out.AttrValidNsec)
	assert.Equal(t, uint64(3600), out.EntryValid)
	assert.Equal(t, uint32(0), out.EntryValidNsec)
	assert.Equal(t, uint64(1), out.Attr.Ino)
	assert.Equal(t, uint64(1000), out.Attr.Size)
}

func TestToFuseAttrOut(t *testing.T) {
	in := &pb.AttrOut{
		AttrValid:     3600,
		AttrValidNsec: 0,
		Attr: &pb.Attr{
			Ino:  1,
			Size: 1000,
			Owner: &pb.Owner{
				Uid: 1000,
				Gid: 1000,
			},
		},
	}

	out := &fuse.AttrOut{}
	toFuseAttrOut(out, in)

	assert.Equal(t, uint64(3600), out.AttrValid)
	assert.Equal(t, uint32(0), out.AttrValidNsec)
	assert.Equal(t, uint64(1), out.Attr.Ino)
	assert.Equal(t, uint64(1000), out.Attr.Size)
}

func TestToFuseOpenOut(t *testing.T) {
	tests := []struct {
		name string
		in   *pb.OpenOut
		want *fuse.OpenOut
	}{
		{
			name: "normal case",
			in: &pb.OpenOut{
				Fh:        123,
				OpenFlags: 1,
				Padding:   0,
			},
			want: &fuse.OpenOut{
				Fh:        123,
				OpenFlags: 1,
				Padding:   0,
			},
		},
		{
			name: "zero values",
			in: &pb.OpenOut{
				Fh:        0,
				OpenFlags: 0,
				Padding:   0,
			},
			want: &fuse.OpenOut{
				Fh:        0,
				OpenFlags: 0,
				Padding:   0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &fuse.OpenOut{}
			toFuseOpenOut(got, tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDealGrpcError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected fuse.Status
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: fuse.OK,
		},
		{
			name:     "unimplemented error",
			err:      status.Error(codes.Unimplemented, "not implemented"),
			expected: fuse.ENOSYS,
		},
		{
			name:     "other gRPC error",
			err:      status.Error(codes.Internal, "internal error"),
			expected: fuse.EIO,
		},
		{
			name:     "non-gRPC error",
			err:      assert.AnError,
			expected: fuse.EIO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dealGrpcError("TestMethod", tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
