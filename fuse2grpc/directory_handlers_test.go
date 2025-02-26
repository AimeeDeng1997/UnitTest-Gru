package fuse2grpc

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestDeltaSize(t *testing.T) {
	testCases := []struct {
		name     string
		dirent   _Dirent
		expected int
	}{
		{
			name: "empty name",
			dirent: _Dirent{
				Ino:     1,
				Off:     0,
				NameLen: 0,
				Typ:     0,
			},
			expected: 12, // 4 (Mode) + 8 (Ino) + 0 (NameLen)
		},
		{
			name: "short name",
			dirent: _Dirent{
				Ino:     123,
				Off:     456,
				NameLen: 5,
				Typ:     8,
			},
			expected: 17, // 4 (Mode) + 8 (Ino) + 5 (NameLen)
		},
		{
			name: "long name",
			dirent: _Dirent{
				Ino:     999999,
				Off:     888888,
				NameLen: 255,
				Typ:     4,
			},
			expected: 267, // 4 (Mode) + 8 (Ino) + 255 (NameLen)
		},
		{
			name: "max name length",
			dirent: _Dirent{
				Ino:     1,
				Off:     1,
				NameLen: 1024,
				Typ:     0,
			},
			expected: 1036, // 4 (Mode) + 8 (Ino) + 1024 (NameLen)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			direntPtr := (*_Dirent)(unsafe.Pointer(&tc.dirent))
			result := deltaSize(direntPtr)
			require.Equal(t, tc.expected, result)
		})
	}
}
