package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	assert.Equal(t,
		"unsafe/fit-in/800x800/filters%3Afill%28white%29%3Awatermark%28raw.githubusercontent.com/cshum/imagor/master/testdata/gopher.png%2Crepeat%2Cbottom%2C10%29%3Aformat%28jpeg%29/https%3A/raw.githubusercontent.com/golang-samples/gopher-vector/master/gopher+.png",
		Normalize("/unsafe/fit-in/800x800/filters:fill(white):watermark(raw.githubusercontent.com/cshum/imagor/master/testdata/gopher.png,repeat,bottom,10):format(jpeg)/https://raw.githubusercontent.com/golang-samples/gopher-vector/master/gopher .png///", nil),
	)

	assert.Equal(t,
		"unsafe/fit-in/800x800/filters%3Afill%28white%29%3Awatermark%28raw.githubusercontent.com/cshum/imagor/master/testdata/gopher.png%2Crepeat%2Cbottom%2C10%29%3Aformat%28jpeg%29/https%3A/raw.githubusercontent.com/golang-samples/gopher-vector/master/gopher .png",
		Normalize("/unsafe/fit-in/800x800/filters:fill(white):watermark(raw.githubusercontent.com/cshum/imagor/master/testdata/gopher.png,repeat,bottom,10):format(jpeg)/https://raw.githubusercontent.com/golang-samples/gopher-vector/master/gopher .png///", NewSafeChars(" ")),
		"should exclude escape space",
	)

	assert.Equal(t, "a+", Normalize("a ", nil))
}
