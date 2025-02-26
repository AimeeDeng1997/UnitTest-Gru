// Copyright 2016 the Go-FUSE Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fuse2grpc

import (
	"os"
	"sync"
	"testing"
)

func TestGetPool(t *testing.T) {
	bp := &bufferPool{}

	// Test initial pool creation
	pool1 := bp.getPool(1)
	if pool1 == nil {
		t.Error("Expected non-nil pool")
	}

	// Test getting same pool again
	pool2 := bp.getPool(1)
	if pool1 != pool2 {
		t.Error("Expected same pool instance")
	}

	// Test getting different sized pool
	pool3 := bp.getPool(2)
	if pool3 == pool1 {
		t.Error("Expected different pool instance")
	}

	// Test concurrent access
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bp.getPool(1)
		}()
	}
	wg.Wait()
}

func TestAllocBuffer(t *testing.T) {
	bp := &bufferPool{}
	pageSize := os.Getpagesize()

	tests := []struct {
		name     string
		size     uint32
		expected int
	}{
		{"Zero size", 0, pageSize},
		{"Small size", 100, pageSize},
		{"Page size", uint32(pageSize), pageSize},
		{"Large size", uint32(pageSize + 100), pageSize * 2},
		{"Multiple pages", uint32(pageSize * 3), pageSize * 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bp.AllocBuffer(tt.size)
			if cap(buf) < tt.expected {
				t.Errorf("AllocBuffer(%d) got capacity %d, want >= %d", tt.size, cap(buf), tt.expected)
			}
			if len(buf) != int(tt.size) {
				t.Errorf("AllocBuffer(%d) got length %d, want %d", tt.size, len(buf), tt.size)
			}
		})
	}
}

func TestFreeBuffer(t *testing.T) {
	bp := &bufferPool{}
	pageSize := os.Getpagesize()

	tests := []struct {
		name string
		buf  []byte
	}{
		{"Nil buffer", nil},
		{"Zero capacity buffer", make([]byte, 0)},
		{"Non-page aligned buffer", make([]byte, 100)},
		{"Page aligned buffer", bp.AllocBuffer(uint32(pageSize))},
		{"Multiple page buffer", bp.AllocBuffer(uint32(pageSize * 2))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify it doesn't panic
			bp.FreeBuffer(tt.buf)
		})
	}

	// Test reuse of freed buffer
	size := uint32(pageSize)
	buf1 := bp.AllocBuffer(size)
	bp.FreeBuffer(buf1)
	buf2 := bp.AllocBuffer(size)
	if cap(buf2) != pageSize {
		t.Errorf("Expected reused buffer capacity %d, got %d", pageSize, cap(buf2))
	}
}

func TestBufferPoolConcurrent(t *testing.T) {
	bp := &bufferPool{}
	pageSize := os.Getpagesize()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buf := bp.AllocBuffer(uint32(pageSize))
			bp.FreeBuffer(buf)
		}()
	}
	wg.Wait()
}
