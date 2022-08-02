package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultipleOptionsSet(t *testing.T) {
	b := NewServerOptionsBuilder()
	b.AddNetwork("tcp").AddAddress("http")

	assert.Len(t, b.srvOpts, 2)
}
