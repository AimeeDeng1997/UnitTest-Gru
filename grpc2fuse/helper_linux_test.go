package grpc2fuse

import (
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"

	"github.com/chiyutianyi/grpcfuse/pb"
)

func TestGetUmask(t *testing.T) {
	tests := []struct {
		name     string
		input    *fuse.MknodIn
		expected uint16
	}{
		{
			name: "normal case",
			input: &fuse.MknodIn{
				Umask: 0022,
			},
			expected: 0022,
		},
		{
			name: "zero umask",
			input: &fuse.MknodIn{
				Umask: 0,
			},
			expected: 0,
		},
		{
			name: "all bits set",
			input: &fuse.MknodIn{
				Umask: 0777,
			},
			expected: 0777,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getUmask(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSetFlags(t *testing.T) {
	// Test the empty function
	out := &fuse.Attr{}
	setFlags(out, 123)
	// No assertions since function is empty
}

func TestSetBlksize(t *testing.T) {
	tests := []struct {
		name     string
		size     uint32
		expected uint32
	}{
		{
			name:     "normal size",
			size:     4096,
			expected: 4096,
		},
		{
			name:     "zero size",
			size:     0,
			expected: 0,
		},
		{
			name:     "large size",
			size:     1048576,
			expected: 1048576,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &fuse.Attr{}
			setBlksize(out, tt.size)
			assert.Equal(t, tt.expected, out.Blksize)
		})
	}
}

func TestSetPadding(t *testing.T) {
	tests := []struct {
		name     string
		padding  uint32
		expected uint32
	}{
		{
			name:     "normal padding",
			padding:  8,
			expected: 8,
		},
		{
			name:     "zero padding",
			padding:  0,
			expected: 0,
		},
		{
			name:     "large padding",
			padding:  1024,
			expected: 1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &fuse.Attr{}
			setPadding(out, tt.padding)
			assert.Equal(t, tt.expected, out.Padding)
		})
	}
}

func TestToPbReadIn(t *testing.T) {
	tests := []struct {
		name     string
		input    *fuse.ReadIn
		expected *pb.ReadIn
	}{
		{
			name: "normal case",
			input: &fuse.ReadIn{
				InHeader: fuse.InHeader{
					Length:  100,
					Opcode: 1,
					Unique: 123,
					NodeId: 456,
				},
				Fh:        789,
				ReadFlags: 0,
				Offset:    1024,
				Size:      4096,
				LockOwner: 111,
				Flags:     222,
				Padding:   333,
			},
			expected: &pb.ReadIn{
				Header: &pb.InHeader{
					Length: 100,
					Opcode: 1,
					Unique: 123,
					NodeId: 456,
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
				Fh:        789,
				ReadFlags: 0,
				Offset:    1024,
				Size:      4096,
				LockOwner: 111,
				Flags:     222,
				Padding:   333,
			},
		},
		{
			name: "zero values",
			input: &fuse.ReadIn{
				InHeader:  fuse.InHeader{},
				Fh:       0,
				Offset:   0,
				Size:     0,
				LockOwner: 0,
				Flags:    0,
				Padding:  0,
			},
			expected: &pb.ReadIn{
				Header: &pb.InHeader{
					Caller: &pb.Caller{
						Owner: &pb.Owner{},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toPbReadIn(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
