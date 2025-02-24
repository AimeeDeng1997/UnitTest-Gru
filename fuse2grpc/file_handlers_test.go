package fuse2grpc

import (
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
)

type mockRawFileSystem struct {
	fuse.RawFileSystem
}

func TestNewServer(t *testing.T) {
	fs := &mockRawFileSystem{}
	s := NewServer(fs)
	if s == nil {
		t.Error("NewServer returned nil")
	}
	if s.fs != fs {
		t.Error("NewServer did not set fs field correctly")
	}
}

func TestSetMsgSizeThreshold(t *testing.T) {
	s := &server{msgSizeThreshold: msgSizeThreshold}
	newThreshold := 2 * msgSizeThreshold
	s.SetMsgSizeThreshold(newThreshold)
	if s.msgSizeThreshold != newThreshold {
		t.Errorf("SetMsgSizeThreshold did not set threshold correctly, got %d want %d",
			s.msgSizeThreshold, newThreshold)
	}
}
