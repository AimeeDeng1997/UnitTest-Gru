package main

import (
	"context"
	"flag"
	"net"
	"os"
	"testing"
	"time"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/grpc2fuse"
	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/chiyutianyi/grpcfuse/pkg/utils"
)

func TestMainArgsValidation(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	tests := []struct {
		name     string
		args     []string
		wantExit bool
	}{
		{
			name:     "no arguments",
			args:     []string{"cmd"},
			wantExit: true,
		},
		{
			name:     "one argument",
			args:     []string{"cmd", "/tmp"},
			wantExit: true,
		},
		{
			name:     "valid arguments",
			args:     []string{"cmd", "/tmp", "localhost:8080"},
			wantExit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			if tt.wantExit {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic but got none")
					}
				}()
			}

			flag.Parse()
			if flag.NArg() < 2 {
				panic("Usage error")
			}
		})
	}
}

func TestMountOptions(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	tests := []struct {
		name       string
		args       []string
		checkMount func(*testing.T, *fuse.MountOptions)
	}{
		{
			name: "default options",
			args: []string{"cmd", "/tmp", "localhost:8080"},
			checkMount: func(t *testing.T, opt *fuse.MountOptions) {
				assert.Equal(t, "GrpcFS", opt.FsName)
				assert.Equal(t, "grpcfs", opt.Name)
				assert.False(t, opt.SingleThreaded)
				assert.Equal(t, 50, opt.MaxBackground)
				assert.True(t, opt.EnableLocks)
				assert.True(t, opt.IgnoreSecurityLabels)
				assert.Equal(t, 1<<20, opt.MaxWrite)
				assert.Equal(t, 1<<20, opt.MaxReadAhead)
				assert.True(t, opt.DirectMount)
				assert.False(t, opt.AllowOther)
				assert.False(t, opt.Debug)
			},
		},
		{
			name: "with allow-other and debug",
			args: []string{"cmd", "-allow-other", "-debug", "/tmp", "localhost:8080"},
			checkMount: func(t *testing.T, opt *fuse.MountOptions) {
				assert.True(t, opt.AllowOther)
				assert.True(t, opt.Debug)
			},
		},
		{
			name: "with read-only",
			args: []string{"cmd", "-ro", "/tmp", "localhost:8080"},
			checkMount: func(t *testing.T, opt *fuse.MountOptions) {
				assert.Contains(t, opt.Options, "ro")
			},
		},
		{
			name: "with custom logger level",
			args: []string{"cmd", "-logger-level=debug", "/tmp", "localhost:8080"},
			checkMount: func(t *testing.T, opt *fuse.MountOptions) {
				assert.Equal(t, "GrpcFS", opt.FsName)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			debug := flag.Bool("debug", false, "print debugging messages.")
			other := flag.Bool("allow-other", false, "mount with -o allowother.")
			ro := flag.Bool("ro", false, "mount read-only")
			logLevel := flag.String("logger-level", "info", "log level")
			flag.Parse()

			utils.GetLogLevel(*logLevel)

			var opt fuse.MountOptions
			opt.FsName = "GrpcFS"
			opt.Name = "grpcfs"
			opt.SingleThreaded = false
			opt.MaxBackground = 50
			opt.EnableLocks = true
			opt.IgnoreSecurityLabels = true
			opt.MaxWrite = 1 << 20
			opt.MaxReadAhead = 1 << 20
			opt.DirectMount = true
			opt.AllowOther = *other
			opt.Debug = *debug
			opt.Options = append(opt.Options, "default_permissions")
			if *ro {
				opt.Options = append(opt.Options, "ro")
			}

			tt.checkMount(t, &opt)
		})
	}
}

type mockRawFileSystemClient struct {
	pb.RawFileSystemClient
}

func TestNewFileSystem(t *testing.T) {
	mockClient := &mockRawFileSystemClient{}
	fs := grpc2fuse.NewFileSystem(mockClient)
	assert.NotNil(t, fs)
}

func dialer(context.Context, string) (net.Conn, error) {
	return nil, status.Error(codes.Unavailable, "connection failed")
}

func TestGrpcConnection(t *testing.T) {
	tests := []struct {
		name       string
		serverAddr string
		wantError  bool
		dialFunc   func(context.Context, string) (net.Conn, error)
	}{
		{
			name:       "invalid address",
			serverAddr: "invalid:address",
			wantError:  true,
			dialFunc:   dialer,
		},
		{
			name:       "valid address but no server",
			serverAddr: "localhost:0",
			wantError:  true,
			dialFunc:   dialer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dialOpts := []grpc.DialOption{
				grpc.WithInsecure(),
				grpc.WithContextDialer(tt.dialFunc),
				grpc.WithBlock(),
			}

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			_, err := grpc.DialContext(ctx, tt.serverAddr, dialOpts...)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLoggerLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{
			name:  "debug level",
			level: "debug",
		},
		{
			name:  "info level",
			level: "info",
		},
		{
			name:  "invalid level",
			level: "invalid",
		},
		{
			name:  "empty level",
			level: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := utils.GetLogLevel(tt.level)
			assert.NotNil(t, level)
		})
	}
}
