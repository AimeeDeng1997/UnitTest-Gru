package grpc2fuse_test

import (
	"errors"
	"testing"

	"github.com/chiyutianyi/grpcfuse/pb"
	"github.com/stretchr/testify/assert"
)

// MockReadDirClient implements RawFileSystem_ReadDirClient interface for testing
type MockReadDirClient struct {
	responses []*pb.ReadDirResponse
	index     int
	err       error
}

func NewMockReadDirClient(responses []*pb.ReadDirResponse, err error) *MockReadDirClient {
	return &MockReadDirClient{
		responses: responses,
		index:     0,
		err:       err,
	}
}

func (m *MockReadDirClient) Recv() (*pb.ReadDirResponse, error) {
	if m.err != nil {
		return nil, m.err
	}

	if m.index >= len(m.responses) {
		return nil, nil
	}

	response := m.responses[m.index]
	m.index++
	return response, nil
}

func TestReadDirClient_Recv(t *testing.T) {
	t.Run("successful recv", func(t *testing.T) {
		responses := []*pb.ReadDirResponse{
			{
				Status: &pb.Status{Code: 0},
				Entries: []*pb.DirEntry{
					{
						Ino:  1,
						Name: []byte("file1"),
						Mode: 0644,
					},
				},
			},
		}

		client := NewMockReadDirClient(responses, nil)
		resp, err := client.Recv()

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, int32(0), resp.Status.Code)
		assert.Len(t, resp.Entries, 1)
		assert.Equal(t, uint64(1), resp.Entries[0].Ino)
		assert.Equal(t, []byte("file1"), resp.Entries[0].Name)
		assert.Equal(t, uint32(0644), resp.Entries[0].Mode)
	})

	t.Run("error on recv", func(t *testing.T) {
		expectedErr := errors.New("mock error")
		client := NewMockReadDirClient(nil, expectedErr)

		resp, err := client.Recv()

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, resp)
	})

	t.Run("empty response", func(t *testing.T) {
		client := NewMockReadDirClient([]*pb.ReadDirResponse{}, nil)

		resp, err := client.Recv()

		assert.NoError(t, err)
		assert.Nil(t, resp)
	})
}
