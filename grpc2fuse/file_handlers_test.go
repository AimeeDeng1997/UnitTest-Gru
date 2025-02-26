package grpc2fuse_test

import (
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hanwen/go-fuse/v2/fuse"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/chiyutianyi/grpcfuse/grpc2fuse"
	"github.com/chiyutianyi/grpcfuse/mock"
	"github.com/chiyutianyi/grpcfuse/pb"
)

var testInHeader = fuse.InHeader{
	Length: 100,
	Opcode: 1,
	Unique: 1,
	NodeId: 1,
}

func TestOpen(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock.NewMockRawFileSystemClient(ctrl)
	fs := grpc2fuse.NewFileSystem(client)
	log.SetLevel(log.ErrorLevel)

	tests := []struct {
		name       string
		response   *pb.OpenResponse
		err        error
		wantStatus fuse.Status
	}{
		{
			name: "success",
			response: &pb.OpenResponse{
				Status:  &pb.Status{Code: 0},
				OpenOut: &pb.OpenOut{Fh: 123, OpenFlags: 0},
			},
			wantStatus: fuse.OK,
		},
		{
			name: "error status",
			response: &pb.OpenResponse{
				Status: &pb.Status{Code: int32(fuse.EACCES)},
			},
			wantStatus: fuse.EACCES,
		},
		{
			name:       "grpc error",
			err:        io.EOF,
			wantStatus: fuse.EIO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &fuse.OpenIn{
				InHeader: testInHeader,
				Flags:    0,
				Mode:     0644,
			}
			out := &fuse.OpenOut{}

			client.EXPECT().Open(gomock.Any(), gomock.Any()).Return(tt.response, tt.err)
			status := fs.Open(nil, in, out)
			require.Equal(t, tt.wantStatus, status)

			if tt.response != nil && tt.response.Status.Code == 0 {
				require.Equal(t, tt.response.OpenOut.Fh, out.Fh)
				require.Equal(t, tt.response.OpenOut.OpenFlags, out.OpenFlags)
			}
		})
	}
}

func TestRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock.NewMockRawFileSystemClient(ctrl)
	fs := grpc2fuse.NewFileSystem(client)
	log.SetLevel(log.ErrorLevel)

	tests := []struct {
		name       string
		messages   []struct {
			response *pb.ReadResponse
			err      error
		}
		wantData   string
		wantStatus fuse.Status
	}{
		{
			name: "success",
			messages: []struct {
				response *pb.ReadResponse
				err      error
			}{
				{&pb.ReadResponse{Status: &pb.Status{Code: 0}, Buffer: []byte("hello ")}, nil},
				{&pb.ReadResponse{Status: &pb.Status{Code: 0}, Buffer: []byte("world")}, nil},
				{nil, io.EOF},
			},
			wantData:   "hello world",
			wantStatus: fuse.OK,
		},
		{
			name: "error status",
			messages: []struct {
				response *pb.ReadResponse
				err      error
			}{
				{&pb.ReadResponse{Status: &pb.Status{Code: int32(fuse.EACCES)}}, nil},
			},
			wantStatus: fuse.EACCES,
		},
		{
			name: "grpc error",
			messages: []struct {
				response *pb.ReadResponse
				err      error
			}{
				{nil, io.ErrUnexpectedEOF},
			},
			wantStatus: fuse.EIO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &fuse.ReadIn{
				InHeader: testInHeader,
				Size:     1000,
			}
			buf := make([]byte, 1000)

			readClient := mock.NewMockRawFileSystem_ReadClient(ctrl)
			client.EXPECT().Read(gomock.Any(), gomock.Any()).Return(readClient, nil)

			var calls []*gomock.Call
			for _, msg := range tt.messages {
				calls = append(calls, readClient.EXPECT().Recv().Return(msg.response, msg.err))
			}
			gomock.InOrder(calls...)

			rs, status := fs.Read(nil, in, buf)
			require.Equal(t, tt.wantStatus, status)

			if status == fuse.OK {
				out, status := rs.Bytes(buf)
				require.Equal(t, fuse.OK, status)
				require.Equal(t, tt.wantData, string(out))
			}
		})
	}
}

func TestLseek(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := mock.NewMockRawFileSystemClient(ctrl)
	fs := grpc2fuse.NewFileSystem(client)
	log.SetLevel(log.ErrorLevel)

	tests := []struct {
		name       string
		response   *pb.LseekResponse
		err        error
		wantStatus fuse.Status
		wantOffset uint64
	}{
		{
			name: "success",
			response: &pb.LseekResponse{
				Status: &pb.Status{Code: 0},
				Offset: 100,
			},
			wantStatus: fuse.OK,
			wantOffset: 100,
		},
		{
			name: "error status",
			response: &pb.LseekResponse{
				Status: &pb.Status{Code: int32(fuse.EINVAL)},
			},
			wantStatus: fuse.EINVAL,
		},
		{
			name:       "grpc error",
			err:        io.EOF,
			wantStatus: fuse.EIO,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &fuse.LseekIn{
				InHeader: testInHeader,
				Offset:   0,
				Whence:   0,
			}
			out := &fuse.LseekOut{}

			client.EXPECT().Lseek(gomock.Any(), gomock.Any()).Return(tt.response, tt.err)
			status := fs.Lseek(nil, in, out)
			require.Equal(t, tt.wantStatus, status)

			if status == fuse.OK {
				require.Equal(t, tt.wantOffset, out.Offset)
			}
		})
	}
}
