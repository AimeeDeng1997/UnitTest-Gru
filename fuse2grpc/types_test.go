package fuse2grpc

import (
	"testing"
	"unsafe"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
)

func TestDirentSize(t *testing.T) {
	d := _Dirent{}
	expectedSize := uint32(unsafe.Sizeof(d))
	assert.Equal(t, expectedSize, direntSize, "Dirent size should match unsafe.Sizeof")
}

func TestEntryOutSize(t *testing.T) {
	e := fuse.EntryOut{}
	expectedSize := uint32(unsafe.Sizeof(e))
	assert.Equal(t, expectedSize, entryOutSize, "EntryOut size should match unsafe.Sizeof")
}

func TestDirentStructure(t *testing.T) {
	d := _Dirent{
		Ino:     123,
		Off:     456,
		NameLen: 10,
		Typ:     8,
	}

	assert.Equal(t, uint64(123), d.Ino, "Ino field should be set correctly")
	assert.Equal(t, uint64(456), d.Off, "Off field should be set correctly")
	assert.Equal(t, uint32(10), d.NameLen, "NameLen field should be set correctly")
	assert.Equal(t, uint32(8), d.Typ, "Typ field should be set correctly")
}

func TestDirEntryList(t *testing.T) {
	buf := make([]byte, 1024)
	list := &DirEntryList{
		buf:    buf,
		size:   1024,
		offset: 100,
	}

	assert.Equal(t, buf, list.buf, "Buffer should be set correctly")
	assert.Equal(t, 1024, list.size, "Size should be set correctly")
	assert.Equal(t, uint64(100), list.offset, "Offset should be set correctly")
	assert.Nil(t, list.lastDirent, "LastDirent should be nil initially")
}

func TestDirEntryListWithDirent(t *testing.T) {
	buf := make([]byte, int(direntSize)+10) // Buffer large enough for dirent and name
	list := &DirEntryList{
		buf:    buf,
		size:   len(buf),
		offset: 0,
	}

	// Create a dirent in the buffer
	dirent := (*_Dirent)(unsafe.Pointer(&buf[0]))
	dirent.Ino = 1
	dirent.Off = 1
	dirent.NameLen = 5
	dirent.Typ = 4

	list.lastDirent = dirent

	assert.NotNil(t, list.lastDirent, "LastDirent should not be nil")
	assert.Equal(t, uint64(1), list.lastDirent.Ino, "Dirent Ino should be set correctly")
	assert.Equal(t, uint64(1), list.lastDirent.Off, "Dirent Off should be set correctly")
	assert.Equal(t, uint32(5), list.lastDirent.NameLen, "Dirent NameLen should be set correctly")
	assert.Equal(t, uint32(4), list.lastDirent.Typ, "Dirent Typ should be set correctly")
}
