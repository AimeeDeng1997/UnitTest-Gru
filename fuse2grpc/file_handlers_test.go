package fuse2grpc_test

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/fuse2grpc"
	"github.com/chiyutianyi/grpcfuse/pb"
)

type mockFS struct {
	mock.Mock
	fuse.RawFileSystem
}

func (m *mockFS) Create(cancel <-chan struct{}, input *fuse.CreateIn, name string, out *fuse.CreateOut) fuse.Status {
	args := m.Called(cancel, input, name, out)
	return args.Get(0).(fuse.Status)
}

func (m *mockFS) Open(cancel <-chan struct{}, input *fuse.OpenIn, out *fuse.OpenOut) fuse.Status {
	args := m.Called(cancel, input, out)
	return args.Get(0).(fuse.Status)
}

func (m *mockFS) Read(cancel <-chan struct{}, input *fuse.ReadIn, buf []byte) (fuse.ReadResult, fuse.Status) {
	args := m.Called(cancel, input, buf)
	if args.Get(0) == nil {
		return nil, args.Get(1).(fuse.Status)
	}
	return args.Get(0).(fuse.ReadResult), args.Get(1).(fuse.Status)
}

func (m *mockFS) Lseek(cancel <-chan struct{}, in *fuse.LseekIn, out *fuse.LseekOut) fuse.Status {
	args := m.Called(cancel, in, out)
	return args.Get(0).(fuse.Status)
}

type mockReadResult struct {
	data []byte
}

func (m *mockReadResult) Bytes(buf []byte) ([]byte, fuse.Status) {
	copy(buf, m.data)
	return m.data, fuse.OK
}

func (m *mockReadResult) Size() int {
	return len(m.data)
}

func (m *mockReadResult) Done() {
}

type mockReadServer struct {
	mock.Mock
	grpc.ServerStream
	ctx context.Context
}

func (m *mockReadServer) Context() context.Context {
	return m.ctx
}

func (m *mockReadServer) Send(resp *pb.ReadResponse) error {
	args := m.Called(resp)
	return args.Error(0)
}

func (m *mockReadServer) RecvMsg(msg interface{}) error {
	return nil
}

func (m *mockReadServer) SendMsg(msg interface{}) error {
	return nil
}

func TestCreate(t *testing.T) {
	mockfs := &mockFS{}
	server := fuse2grpc.NewServer(mockfs)

	tests := []struct {
		name    string
		req     *pb.CreateRequest
		setup   func()
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "success",
			req: &pb.CreateRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
						Pid: 1234,
					},
				},
				Name:  "test.txt",
				Flags: 0,
				Mode:  0644,
			},
			setup: func() {
				mockfs.On("Create", mock.Anything, mock.Anything, "test.txt", mock.Anything).Return(fuse.OK)
			},
			wantErr: false,
		},
		{
			name: "not implemented",
			req: &pb.CreateRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Name: "test.txt",
			},
			setup: func() {
				mockfs.On("Create", mock.Anything, mock.Anything, "test.txt", mock.Anything).Return(fuse.ENOSYS)
			},
			wantErr: true,
			errCode: codes.Unimplemented,
		},
		{
			name: "error",
			req: &pb.CreateRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Name: "test.txt",
			},
			setup: func() {
				mockfs.On("Create", mock.Anything, mock.Anything, "test.txt", mock.Anything).Return(fuse.EACCES)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockfs.ExpectedCalls = nil
			tt.setup()
			resp, err := server.Create(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.errCode, st.Code())
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, resp)
		})
	}
}

func TestOpen(t *testing.T) {
	mockfs := &mockFS{}
	server := fuse2grpc.NewServer(mockfs)

	tests := []struct {
		name    string
		req     *pb.OpenRequest
		setup   func()
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "success",
			req: &pb.OpenRequest{
				OpenIn: &pb.OpenIn{
					Header: &pb.InHeader{
						NodeId: 1,
						Caller: &pb.Caller{
							Owner: &pb.Owner{
								Uid: 1000,
								Gid: 1000,
							},
						},
					},
					Flags: 0,
					Mode:  0644,
				},
			},
			setup: func() {
				mockfs.On("Open", mock.Anything, mock.Anything, mock.Anything).Return(fuse.OK)
			},
			wantErr: false,
		},
		{
			name: "not implemented",
			req: &pb.OpenRequest{
				OpenIn: &pb.OpenIn{
					Header: &pb.InHeader{
						NodeId: 1,
						Caller: &pb.Caller{
							Owner: &pb.Owner{
								Uid: 1000,
								Gid: 1000,
							},
						},
					},
				},
			},
			setup: func() {
				mockfs.On("Open", mock.Anything, mock.Anything, mock.Anything).Return(fuse.ENOSYS)
			},
			wantErr: true,
			errCode: codes.Unimplemented,
		},
		{
			name: "error",
			req: &pb.OpenRequest{
				OpenIn: &pb.OpenIn{
					Header: &pb.InHeader{
						NodeId: 1,
						Caller: &pb.Caller{
							Owner: &pb.Owner{
								Uid: 1000,
								Gid: 1000,
							},
						},
					},
				},
			},
			setup: func() {
				mockfs.On("Open", mock.Anything, mock.Anything, mock.Anything).Return(fuse.EACCES)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockfs.ExpectedCalls = nil
			tt.setup()
			resp, err := server.Open(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.errCode, st.Code())
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, resp)
		})
	}
}

func TestRead(t *testing.T) {
	mockfs := &mockFS{}
	server := fuse2grpc.NewServer(mockfs)

	tests := []struct {
		name    string
		req     *pb.ReadRequest
		setup   func()
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "success",
			req: &pb.ReadRequest{
				ReadIn: &pb.ReadIn{
					Header: &pb.InHeader{
						NodeId: 1,
						Caller: &pb.Caller{
							Owner: &pb.Owner{
								Uid: 1000,
								Gid: 1000,
							},
						},
					},
					Size: 1024,
				},
			},
			setup: func() {
				data := []byte("test data")
				mockfs.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(&mockReadResult{data: data}, fuse.OK)
			},
			wantErr: false,
		},
		{
			name: "not implemented",
			req: &pb.ReadRequest{
				ReadIn: &pb.ReadIn{
					Header: &pb.InHeader{
						NodeId: 1,
						Caller: &pb.Caller{
							Owner: &pb.Owner{
								Uid: 1000,
								Gid: 1000,
							},
						},
					},
				},
			},
			setup: func() {
				mockfs.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(nil, fuse.ENOSYS)
			},
			wantErr: true,
			errCode: codes.Unimplemented,
		},
		{
			name: "read error",
			req: &pb.ReadRequest{
				ReadIn: &pb.ReadIn{
					Header: &pb.InHeader{
						NodeId: 1,
						Caller: &pb.Caller{
							Owner: &pb.Owner{
								Uid: 1000,
								Gid: 1000,
							},
						},
					},
				},
			},
			setup: func() {
				mockfs.On("Read", mock.Anything, mock.Anything, mock.Anything).Return(nil, fuse.EIO)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockfs.ExpectedCalls = nil
			tt.setup()

			mockStream := &mockReadServer{ctx: context.Background()}
			mockStream.On("Send", mock.Anything).Return(nil)

			err := server.Read(tt.req, mockStream)
			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.errCode, st.Code())
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestLseek(t *testing.T) {
	mockfs := &mockFS{}
	server := fuse2grpc.NewServer(mockfs)

	tests := []struct {
		name    string
		req     *pb.LseekRequest
		setup   func()
		wantErr bool
		errCode codes.Code
	}{
		{
			name: "success",
			req: &pb.LseekRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
				Offset: 100,
				Whence: 0,
			},
			setup: func() {
				mockfs.On("Lseek", mock.Anything, mock.Anything, mock.Anything).Return(fuse.OK)
			},
			wantErr: false,
		},
		{
			name: "not implemented",
			req: &pb.LseekRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
			},
			setup: func() {
				mockfs.On("Lseek", mock.Anything, mock.Anything, mock.Anything).Return(fuse.ENOSYS)
			},
			wantErr: true,
			errCode: codes.Unimplemented,
		},
		{
			name: "error",
			req: &pb.LseekRequest{
				Header: &pb.InHeader{
					NodeId: 1,
					Caller: &pb.Caller{
						Owner: &pb.Owner{
							Uid: 1000,
							Gid: 1000,
						},
					},
				},
			},
			setup: func() {
				mockfs.On("Lseek", mock.Anything, mock.Anything, mock.Anything).Return(fuse.EINVAL)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockfs.ExpectedCalls = nil
			tt.setup()
			resp, err := server.Lseek(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.errCode, st.Code())
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, resp)
		})
	}
}
