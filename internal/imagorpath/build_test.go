package imagorpath_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cshum/imagor/internal/imagorpath"
	"github.com/cshum/imagor/pkg/params"
)

type MockSigner struct{}

func (s MockSigner) Sign(_ string) string {
	return "SIGNED"
}

func TestBuildUnprefixed(t *testing.T) {
	p := params.Params{
		Meta:          true,
		Trim:          true,
		TrimBy:        params.TrimByTopLeft,
		TrimTolerance: 5,
		CropLeft:      10.5,
		CropTop:       20.5,
		CropRight:     30.5,
		CropBottom:    40.5,
		FitIn:         true,
		Stretch:       false,
		Width:         800,
		Height:        600,
		PaddingLeft:   10,
		PaddingTop:    20,
		PaddingRight:  10,
		PaddingBottom: 20,
		HFlip:         false,
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

	expected := "meta/trim/top-left:5/10.5x20.5:30.5x40.5/fit-in/800x-600/10x20:10x20/left/bottom/smart/filters:filter1(arg1):filter2(arg2)/example.jpg"
	actual := imagorpath.BuildUnprefixed(p)
	assert.Equal(t, expected, actual, "BuildUnprefixed should generate the correct path")
}

func TestBuildUnsafe(t *testing.T) {
	p := params.Params{
		Meta:          true,
		Trim:          true,
		TrimBy:        params.TrimByTopLeft,
		TrimTolerance: 5,
		CropLeft:      10.5,
		CropTop:       20.5,
		CropRight:     30.5,
		CropBottom:    40.5,
		FitIn:         true,
		Stretch:       false,
		Width:         800,
		Height:        600,
		PaddingLeft:   10,
		PaddingTop:    20,
		PaddingRight:  10,
		PaddingBottom: 20,
		HFlip:         false,
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

	expected := "unsafe/meta/trim/top-left:5/10.5x20.5:30.5x40.5/fit-in/800x-600/10x20:10x20/left/bottom/smart/filters:filter1(arg1):filter2(arg2)/example.jpg"
	actual := imagorpath.BuildUnsafe(p)
	assert.Equal(t, expected, actual, "BuildUnsafe should generate the correct unsafe path")
}

func TestBuildSigned(t *testing.T) {
	p := params.Params{
		Meta:          true,
		Trim:          true,
		TrimBy:        params.TrimByTopLeft,
		TrimTolerance: 5,
		CropLeft:      10.5,
		CropTop:       20.5,
		CropRight:     30.5,
		CropBottom:    40.5,
		FitIn:         true,
		Stretch:       false,
		Width:         800,
		Height:        600,
		PaddingLeft:   10,
		PaddingTop:    20,
		PaddingRight:  10,
		PaddingBottom: 20,
		HFlip:         false,
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

	signer := MockSigner{}

	expected := "SIGNED/meta/trim/top-left:5/10.5x20.5:30.5x40.5/fit-in/800x-600/10x20:10x20/left/bottom/smart/filters:filter1(arg1):filter2(arg2)/example.jpg"
	actual := imagorpath.BuildSigned(p, signer)
	assert.Equal(t, expected, actual)
}
