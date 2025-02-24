package fuse2grpc_test

import "github.com/hanwen/go-fuse/v2/fuse"

type rawFileSystem struct {
	fuse.RawFileSystem
}

func (fs *rawFileSystem) GetXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string, dest []byte) (uint32, fuse.Status) {
	return 0, fuse.ENOSYS
}

func (fs *rawFileSystem) ListXAttr(cancel <-chan struct{}, header *fuse.InHeader, dest []byte) (uint32, fuse.Status) {
	return 0, fuse.ENOSYS
}

func (fs *rawFileSystem) RemoveXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string) fuse.Status {
	return fuse.ENOSYS
}

func (fs *rawFileSystem) String() string {
	return "rawFileSystem"
}
