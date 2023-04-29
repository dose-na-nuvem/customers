package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	// test
	msg := message()

	// verify
	assert.Equal(t, "hello world", msg)
}
