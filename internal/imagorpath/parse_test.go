package imagorpath_test

import (
	"testing"

	"github.com/cshum/imagor/internal/imagorpath"
	"github.com/cshum/imagor/pkg/params"
	"github.com/stretchr/testify/assert"
)

func TestParseUnsafe(t *testing.T) {
	path := "unsafe/meta/trim/top-left:5/10.5x20.5:30.5x40.5/fit-in/stretch/-800x-600/10x20:10x20/left/bottom/smart/filters:filter1(arg1):filter2(arg2)/example.jpg"

	expected := params.Params{
		Unsafe:        true,
		Meta:          true,
		Trim:          true,
		TrimBy:        params.TrimByTopLeft,
		TrimTolerance: 5,
		CropLeft:      10.5,
		CropTop:       20.5,
		CropRight:     30.5,
		CropBottom:    40.5,
		FitIn:         true,
		Stretch:       true,
		Width:         800,
		Height:        600,
		PaddingLeft:   10,
		PaddingTop:    20,
		PaddingRight:  10,
		PaddingBottom: 20,
		HFlip:         true,
		VFlip:         true,
		HAlign:        params.HAlignLeft,
		VAlign:        params.VAlignBottom,
		Smart:         true,
		Filters: params.Filters{
			{Name: "filter1", Args: "arg1"},
			{Name: "filter2", Args: "arg2"},
		},
		Image: "example.jpg",
	}
	actual, err := imagorpath.Parse(path)

	assert.Nil(t, err, "Parse should not return an error")
	assert.Equal(t, expected, actual, "BuildUnprefixed should generate the correct path")
}
