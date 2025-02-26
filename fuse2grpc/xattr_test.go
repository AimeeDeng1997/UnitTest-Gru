package fuse2grpc_test

import (
	"context"
	"testing"

	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chiyutianyi/grpcfuse/fuse2grpc"
	"github.com/chiyutianyi/grpcfuse/pb"
)

type mockFS struct {
	mock.Mock
	fuse.RawFileSystem
}

func (m *mockFS) GetXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string, dest []byte) (uint32, fuse.Status) {
	args := m.Called(cancel, header, attr, dest)
	return args.Get(0).(uint32), args.Get(1).(fuse.Status)
}

func (m *mockFS) ListXAttr(cancel <-chan struct{}, header *fuse.InHeader, dest []byte) (uint32, fuse.Status) {
	args := m.Called(cancel, header, dest)
	return args.Get(0).(uint32), args.Get(1).(fuse.Status)
}

func (m *mockFS) RemoveXAttr(cancel <-chan struct{}, header *fuse.InHeader, attr string) fuse.Status {
	args := m.Called(cancel, header, attr)
	return args.Get(0).(fuse.Status)
}

func (m *mockFS) String() string {
	return "mockFS"
}

func TestGetXAttr(t *testing.T) {
	mockfs := &mockFS{}
	server := fuse2grpc.NewServer(mockfs)

	tests := []struct {
		name     string
		req      *pb.GetXAttrRequest
		mockResp []interface{}
		want     *pb.GetXAttrResponse
		wantErr  error
	}{
		{
			name: "success",
			req: &pb.GetXAttrRequest{
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
				Attr: "user.test",
				Dest: []byte("test"),
			},
			mockResp: []interface{}{uint32(4), fuse.OK},
			want: &pb.GetXAttrResponse{
				Size:   4,
				Status: &pb.Status{Code: 0},
			},
		},
		{
			name: "not implemented",
			req: &pb.GetXAttrRequest{
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
				Attr: "user.test",
			},
			mockResp: []interface{}{uint32(0), fuse.ENOSYS},
			wantErr:  status.Errorf(codes.Unimplemented, "method GetXAttr not implemented"),
		},
		{
			name: "error",
			req: &pb.GetXAttrRequest{
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
				Attr: "user.test",
			},
			mockResp: []interface{}{uint32(0), fuse.ENOENT},
			want: &pb.GetXAttrResponse{
				Size:   0,
				Status: &pb.Status{Code: int32(fuse.ENOENT)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockfs.On("GetXAttr", mock.Anything, mock.Anything, tt.req.Attr, tt.req.Dest).Return(tt.mockResp...).Once()

			got, err := server.GetXAttr(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestListXAttr(t *testing.T) {
	mockfs := &mockFS{}
	server := fuse2grpc.NewServer(mockfs)

	tests := []struct {
		name     string
		req      *pb.ListXAttrRequest
		mockResp []interface{}
		want     *pb.ListXAttrResponse
		wantErr  error
	}{
		{
			name: "success",
			req: &pb.ListXAttrRequest{
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
				Dest: []byte("test"),
			},
			mockResp: []interface{}{uint32(4), fuse.OK},
			want: &pb.ListXAttrResponse{
				Size:   4,
				Status: &pb.Status{Code: 0},
			},
		},
		{
			name: "not implemented",
			req: &pb.ListXAttrRequest{
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
			},
			mockResp: []interface{}{uint32(0), fuse.ENOSYS},
			wantErr:  status.Errorf(codes.Unimplemented, "method ListXAttr not implemented"),
		},
		{
			name: "error",
			req: &pb.ListXAttrRequest{
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
			},
			mockResp: []interface{}{uint32(0), fuse.ENOENT},
			want: &pb.ListXAttrResponse{
				Size:   0,
				Status: &pb.Status{Code: int32(fuse.ENOENT)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockfs.On("ListXAttr", mock.Anything, mock.Anything, tt.req.Dest).Return(tt.mockResp...).Once()

			got, err := server.ListXAttr(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRemoveXAttr(t *testing.T) {
	mockfs := &mockFS{}
	server := fuse2grpc.NewServer(mockfs)

	tests := []struct {
		name     string
		req      *pb.RemoveXAttrRequest
		mockResp []interface{}
		want     *pb.RemoveXAttrResponse
		wantErr  error
	}{
		{
			name: "success",
			req: &pb.RemoveXAttrRequest{
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
				Attr: "user.test",
			},
			mockResp: []interface{}{fuse.OK},
			want: &pb.RemoveXAttrResponse{
				Status: &pb.Status{Code: 0},
			},
		},
		{
			name: "not implemented",
			req: &pb.RemoveXAttrRequest{
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
				Attr: "user.test",
			},
			mockResp: []interface{}{fuse.ENOSYS},
			wantErr:  status.Errorf(codes.Unimplemented, "method RemoveXAttr not implemented"),
		},
		{
			name: "error",
			req: &pb.RemoveXAttrRequest{
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
				Attr: "user.test",
			},
			mockResp: []interface{}{fuse.ENOENT},
			want: &pb.RemoveXAttrResponse{
				Status: &pb.Status{Code: int32(fuse.ENOENT)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockfs.On("RemoveXAttr", mock.Anything, mock.Anything, tt.req.Attr).Return(tt.mockResp...).Once()

			got, err := server.RemoveXAttr(context.Background(), tt.req)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
