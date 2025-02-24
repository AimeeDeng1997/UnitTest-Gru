package main

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	"github.com/chiyutianyi/grpcfuse/fuse2grpc"
	"github.com/chiyutianyi/grpcfuse/pb"
)

func startTestServer(t *testing.T, tmpDir string) (*grpc.Server, net.Listener) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)

	root, err := fs.NewLoopbackRoot(tmpDir)
	assert.NoError(t, err)

	sec := time.Second
	opts := &fs.Options{
		AttrTimeout:  &sec,
		EntryTimeout: &sec,
	}

	rawFS := fs.NewNodeFS(root, opts)
	srv := fuse2grpc.NewServer(rawFS)

	s := grpc.NewServer()
	pb.RegisterRawFileSystemServer(s, srv)

	go s.Serve(l)

	return s, l
}

func TestServerStartup(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "grpcfuse-server-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	s, l := startTestServer(t, tmpDir)
	defer s.Stop()
	defer l.Close()

	conn, err := grpc.Dial(l.Addr().String(), grpc.WithInsecure())
	assert.NoError(t, err)
	defer conn.Close()

	client := pb.NewRawFileSystemClient(conn)

	ctx := context.Background()
	resp, err := client.String(ctx, &pb.StringRequest{})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Value)
}

func TestServerShutdown(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "grpcfuse-shutdown-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	s, l := startTestServer(t, tmpDir)
	defer l.Close()

	time.Sleep(100 * time.Millisecond)
	s.GracefulStop()

	_, err = net.Dial("tcp", l.Addr().String())
	assert.Error(t, err)
}
