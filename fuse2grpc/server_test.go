package fuse2grpc

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"

	"github.com/chiyutianyi/grpcfuse/pb"
)

type mockRawFS struct {
	fuse.RawFileSystem
	stringVal string
}

func (m *mockRawFS) String() string {
	return m.stringVal
}

func TestNewServer(t *testing.T) {
	fs := &mockRawFS{}
	s := NewServer(fs)

	assert.NotNil(t, s)
	assert.Equal(t, fs, s.fs)
	assert.Equal(t, msgSizeThreshold, s.msgSizeThreshold)
}

func TestSetMsgSizeThreshold(t *testing.T) {
	s := &server{msgSizeThreshold: msgSizeThreshold}
	newThreshold := 2 << 20

	s.SetMsgSizeThreshold(newThreshold)
	assert.Equal(t, newThreshold, s.msgSizeThreshold)
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		fsString string
		wantErr  bool
	}{
		{
			name:     "normal string",
			fsString: "test filesystem",
			wantErr:  false,
		},
		{
			name:     "empty string",
			fsString: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &mockRawFS{stringVal: tt.fsString}
			s := NewServer(fs)

			resp, err := s.String(context.Background(), &pb.StringRequest{})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.fsString, resp.Value)
		})
	}
}
