package grpc2fuse

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	cancel := make(chan struct{})
	ctx := newContext(cancel)

	assert.NotNil(t, ctx)
	assert.IsType(t, make(<-chan struct{}), ctx.cancel)
	assert.Equal(t, context.Background(), ctx.Context)
}

func TestFuseContextDone(t *testing.T) {
	cancel := make(chan struct{})
	ctx := newContext(cancel)

	done := ctx.Done()
	assert.IsType(t, make(<-chan struct{}), done)
	assert.Equal(t, ctx.cancel, done)
}

func TestFuseContextErr(t *testing.T) {
	cancel := make(chan struct{})
	ctx := newContext(cancel)

	err := ctx.Err()
	assert.Equal(t, context.Canceled, err)
}
