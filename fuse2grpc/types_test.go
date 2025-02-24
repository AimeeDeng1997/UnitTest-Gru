package fuse2grpc

import (
	"testing"
	"unsafe"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
)

func TestDirentSizeConstant(t *testing.T) {
	// Test direntSize constant matches actual size of _Dirent struct
	var d _Dirent
	expectedSize := unsafe.Sizeof(d)
	assert.Equal(t, uint32(expectedSize), direntSize, "direntSize constant should match actual _Dirent struct size")
}

func TestEntryOutSizeConstant(t *testing.T) {
	// Test entryOutSize constant matches actual size of EntryOut struct
	var e fuse.EntryOut
	expectedSize := unsafe.Sizeof(e)
	assert.Equal(t, uint32(expectedSize), entryOutSize, "entryOutSize constant should match actual EntryOut struct size")
}

func TestDirEntryListStructure(t *testing.T) {
	// Test DirEntryList struct has expected fields
	var del DirEntryList

	// Verify field types
	assert.IsType(t, []byte{}, del.buf, "buf field should be []byte")
	assert.IsType(t, 0, del.size, "size field should be int")
	assert.IsType(t, uint64(0), del.offset, "offset field should be uint64")
	assert.IsType(t, (*_Dirent)(nil), del.lastDirent, "lastDirent field should be pointer to _Dirent struct")
}

func TestDirentFields(t *testing.T) {
	d := _Dirent{
		Ino: 1234,
		Off: 5678,
		NameLen: 10,
		Typ: 1,
	}

	assert.Equal(t, uint64(1234), d.Ino)
	assert.Equal(t, uint64(5678), d.Off)
	assert.Equal(t, uint32(10), d.NameLen)
	assert.Equal(t, uint32(1), d.Typ)
}

func TestDirEntryListInit(t *testing.T) {
	del := DirEntryList{
		buf: make([]byte, 100),
		size: 100,
		offset: 42,
		lastDirent: &_Dirent{},
	}

	assert.Equal(t, 100, len(del.buf))
	assert.Equal(t, 100, del.size)
	assert.Equal(t, uint64(42), del.offset)
	assert.NotNil(t, del.lastDirent)
}
