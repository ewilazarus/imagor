package signer

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSigner(t *testing.T) {
	signer := New("abcd", sha256.New, 28)
	assert.Equal(t, signer.Sign("assfasf"), "zb6uWXQxwJDOe_zOgxkuj96Etrsz")
}
