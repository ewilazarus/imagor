package path

import (
	"crypto/sha256"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/cshum/imagor/internal/params"
	"github.com/cshum/imagor/pkg/signer"
	"github.com/stretchr/testify/assert"
)

func TestParseGenerate(t *testing.T) {
	tests := []struct {
		name   string
		uri    string
		params params.Params
		signer signer.Signer
	}{
		{
			name: "non url image",
			uri:  "meta/trim/10x11:12x13/fit-in/-300x-200/left/top/smart/filters:some_filter()/img",
			params: params.Params{
				Path:       "meta/trim/10x11:12x13/fit-in/-300x-200/left/top/smart/filters:some_filter()/img",
				Image:      "img",
				Trim:       true,
				TrimBy:     "top-left",
				CropLeft:   10,
				CropTop:    11,
				CropRight:  12,
				CropBottom: 13,
				Width:      300,
				Height:     200,
				Meta:       true,
				HFlip:      true,
				VFlip:      true,
				HAlign:     "left",
				VAlign:     "top",
				Smart:      true,
				FitIn:      true,
				Filters:    []params.Filter{{Name: "some_filter"}},
			},
		},
		{
			name: "url image",
			uri:  "meta/trim:bottom-right:100/10x11:12x13/fit-in/-300x-200/left/top/smart/filters:some_filter()/s.glbimg.com/es/ge/f/original/2011/03/29/orlandosilva_60.jpg",
			params: params.Params{
				Path:          "meta/trim:bottom-right:100/10x11:12x13/fit-in/-300x-200/left/top/smart/filters:some_filter()/s.glbimg.com/es/ge/f/original/2011/03/29/orlandosilva_60.jpg",
				Image:         "s.glbimg.com/es/ge/f/original/2011/03/29/orlandosilva_60.jpg",
				Trim:          true,
				TrimBy:        params.TrimByBottomRight,
				TrimTolerance: 100,
				CropLeft:      10,
				CropTop:       11,
				CropRight:     12,
				CropBottom:    13,
				Width:         300,
				Height:        200,
				Meta:          true,
				HFlip:         true,
				VFlip:         true,
				HAlign:        "left",
				VAlign:        "top",
				Smart:         true,
				FitIn:         true,
				Filters:       []params.Filter{{Name: "some_filter"}},
			},
		},
		{
			name: "url in filter",
			uri:  "filters:watermark(s.glbimg.com/es/ge/f/original/2011/03/29/orlandosilva_60.jpg,0,0,0)/img",
			params: params.Params{
				Path:    "filters:watermark(s.glbimg.com/es/ge/f/original/2011/03/29/orlandosilva_60.jpg,0,0,0)/img",
				Image:   "img",
				Filters: []params.Filter{{Name: "watermark", Args: "s.glbimg.com/es/ge/f/original/2011/03/29/orlandosilva_60.jpg,0,0,0"}},
			},
		},
		{
			name: "multiple filters",
			uri:  "filters:watermark(s.glbimg.com/es/ge/f/original/2011/03/29/orlandosilva_60.jpg,0,0,0):brightness(-50):grayscale()/img",
			params: params.Params{
				Path:  "filters:watermark(s.glbimg.com/es/ge/f/original/2011/03/29/orlandosilva_60.jpg,0,0,0):brightness(-50):grayscale()/img",
				Image: "img",
				Filters: []params.Filter{
					{
						Name: "watermark",
						Args: "s.glbimg.com/es/ge/f/original/2011/03/29/orlandosilva_60.jpg,0,0,0",
					},
					{
						Name: "brightness",
						Args: "-50",
					},
					{
						Name: "grayscale",
					},
				},
			},
		},
		{
			name: "nested filters",
			uri:  "filters:watermark(s.glbimg.com/filters:label(abc):watermark(aaa.com/fit-in/filters:aaa(bbb))/aaa.jpg,0,0,0):brightness(-50):grayscale()/img",
			params: params.Params{
				Path:  "filters:watermark(s.glbimg.com/filters:label(abc):watermark(aaa.com/fit-in/filters:aaa(bbb))/aaa.jpg,0,0,0):brightness(-50):grayscale()/img",
				Image: "img",
				Filters: []params.Filter{
					{
						Name: "watermark",
						Args: "s.glbimg.com/filters:label(abc):watermark(aaa.com/fit-in/filters:aaa(bbb))/aaa.jpg,0,0,0",
					},
					{
						Name: "brightness",
						Args: "-50",
					},
					{
						Name: "grayscale",
					},
				},
			},
		},
		{
			name: "filters with unicode",
			uri:  "filters:label(哈哈,1,2,3):brightness(-50):grayscale()/img",
			params: params.Params{
				Path:  "filters:label(哈哈,1,2,3):brightness(-50):grayscale()/img",
				Image: "img",
				Filters: []params.Filter{
					{
						Name: "label",
						Args: "哈哈,1,2,3",
					},
					{
						Name: "brightness",
						Args: "-50",
					},
					{
						Name: "grayscale",
					},
				},
			},
		},

		{
			name: "no params",
			uri:  "unsafe/https://foobar/en/latest/_images/man_before_sharpen.png",
			params: params.Params{
				Path:   "https://foobar/en/latest/_images/man_before_sharpen.png",
				Image:  "https://foobar/en/latest/_images/man_before_sharpen.png",
				Unsafe: true,
			},
		},
		{
			name: "contains query",
			uri:  "unsafe/https%3A%2F%2Ffoobar%2Fen%2Flatest%2F_images%2Fman_before_sharpen.png%3Ffoo%3Dbar",
			params: params.Params{
				Path:   "https%3A%2F%2Ffoobar%2Fen%2Flatest%2F_images%2Fman_before_sharpen.png%3Ffoo%3Dbar",
				Image:  "https://foobar/en/latest/_images/man_before_sharpen.png?foo=bar",
				Unsafe: true,
			},
		},
		{
			name: "image contains keyword trim",
			uri:  "unsafe/trim%2Fimg",
			params: params.Params{
				Path:   "trim%2Fimg",
				Image:  "trim/img",
				Unsafe: true,
			},
		},
		{
			name: "image contains keyword meta",
			uri:  "unsafe/meta%2Fimg",
			params: params.Params{
				Path:   "meta%2Fimg",
				Image:  "meta/img",
				Unsafe: true,
			},
		},
		{
			name: "image contains keyword center",
			uri:  "unsafe/center%2Fimg",
			params: params.Params{
				Path:   "center%2Fimg",
				Image:  "center/img",
				Unsafe: true,
			},
		},
		{
			name: "image contains keyword smart",
			uri:  "unsafe/smart%2Fimg",
			params: params.Params{
				Path:   "smart%2Fimg",
				Image:  "smart/img",
				Unsafe: true,
			},
		},
		{
			name: "image contains keyword fit-in",
			uri:  "unsafe/fit-in%2Fimg",
			params: params.Params{
				Path:   "fit-in%2Fimg",
				Image:  "fit-in/img",
				Unsafe: true,
			},
		},
		{
			name: "image contains keyword stretch",
			uri:  "unsafe/stretch%2Fimg",
			params: params.Params{
				Path:   "stretch%2Fimg",
				Image:  "stretch/img",
				Unsafe: true,
			},
		},
		{
			name: "image contains keyword top",
			uri:  "unsafe/top%2Fimg",
			params: params.Params{
				Path:   "top%2Fimg",
				Image:  "top/img",
				Unsafe: true,
			},
		},
		{
			name: "image contains keyword left",
			uri:  "unsafe/left%2Fimg",
			params: params.Params{
				Path:   "left%2Fimg",
				Image:  "left/img",
				Unsafe: true,
			},
		},
		{
			name: "image contains keyword right",
			uri:  "unsafe/right%2Fimg",
			params: params.Params{
				Path:   "right%2Fimg",
				Image:  "right/img",
				Unsafe: true,
			},
		},
		{
			name: "image contains keyword bottom",
			uri:  "unsafe/bottom%2Fimg",
			params: params.Params{
				Path:   "bottom%2Fimg",
				Image:  "bottom/img",
				Unsafe: true,
			},
		},
		{
			name: "padding without dimensions",
			uri:  "unsafe/fit-in/0x0/5x6:7x8/https://foobar/en/latest/_images/man_before_sharpen.png",
			params: params.Params{
				Path:          "fit-in/0x0/5x6:7x8/https://foobar/en/latest/_images/man_before_sharpen.png",
				Image:         "https://foobar/en/latest/_images/man_before_sharpen.png",
				Unsafe:        true,
				FitIn:         true,
				PaddingLeft:   5,
				PaddingTop:    6,
				PaddingRight:  7,
				PaddingBottom: 8,
			},
		},
		{
			name: "url in filters",
			uri:  "unsafe/stretch/500x350/filters:watermark(http://thumborize.me/static/img/beach.jpg,100,100,50)/http://thumborize.me/static/img/beach.jpg",
			params: params.Params{
				Path:    "stretch/500x350/filters:watermark(http://thumborize.me/static/img/beach.jpg,100,100,50)/http://thumborize.me/static/img/beach.jpg",
				Image:   "http://thumborize.me/static/img/beach.jpg",
				Width:   500,
				Height:  350,
				Unsafe:  true,
				Stretch: true,
				Filters: []params.Filter{
					{
						Name: "watermark",
						Args: "http://thumborize.me/static/img/beach.jpg,100,100,50",
					},
				},
			},
		},
		{
			name:   "non url image with hash",
			uri:    "VTAq7YIRbEXgtwAcsTMhAjvBuT8=/meta/10x11:12x13/fit-in/-300x-200/5x6/left/top/smart/filters:some_filter()/img",
			signer: signer.NewDefaultSigner("1234"),
			params: params.Params{
				Path:          "meta/10x11:12x13/fit-in/-300x-200/5x6/left/top/smart/filters:some_filter()/img",
				Hash:          "VTAq7YIRbEXgtwAcsTMhAjvBuT8=",
				Image:         "img",
				CropLeft:      10,
				CropTop:       11,
				CropRight:     12,
				CropBottom:    13,
				Width:         300,
				Height:        200,
				Meta:          true,
				HFlip:         true,
				VFlip:         true,
				HAlign:        "left",
				VAlign:        "top",
				Smart:         true,
				FitIn:         true,
				PaddingLeft:   5,
				PaddingTop:    6,
				PaddingRight:  5,
				PaddingBottom: 6,
				Filters:       []params.Filter{{Name: "some_filter"}},
			},
		},
		{
			name:   "non url image with hash and custom signer",
			uri:    "XBCO7esuLsNQuSF2v9ie36pESRGx2rzLjhUxXWnV/meta/10x11:12x13/fit-in/-300x-200/5x6/left/top/smart/filters:some_filter()/img",
			signer: signer.NewHMACSigner(sha256.New, 40, "1234"),
			params: params.Params{
				Path:          "meta/10x11:12x13/fit-in/-300x-200/5x6/left/top/smart/filters:some_filter()/img",
				Hash:          "XBCO7esuLsNQuSF2v9ie36pESRGx2rzLjhUxXWnV",
				Image:         "img",
				CropLeft:      10,
				CropTop:       11,
				CropRight:     12,
				CropBottom:    13,
				Width:         300,
				Height:        200,
				Meta:          true,
				HFlip:         true,
				VFlip:         true,
				HAlign:        "left",
				VAlign:        "top",
				Smart:         true,
				FitIn:         true,
				PaddingLeft:   5,
				PaddingTop:    6,
				PaddingRight:  5,
				PaddingBottom: 6,
				Filters:       []params.Filter{{Name: "some_filter"}},
			},
		},
		{
			name: "non url image with crop by percentage",
			uri:  "meta/trim/0.2x0.15:0.45x0.67/fit-in/-300x-200/left/top/smart/filters:some_filter()/img",
			params: params.Params{
				Path:       "meta/trim/0.2x0.15:0.45x0.67/fit-in/-300x-200/left/top/smart/filters:some_filter()/img",
				Image:      "img",
				Trim:       true,
				TrimBy:     "top-left",
				CropLeft:   0.2,
				CropTop:    0.15,
				CropRight:  0.45,
				CropBottom: 0.67,
				Width:      300,
				Height:     200,
				Meta:       true,
				HFlip:      true,
				VFlip:      true,
				HAlign:     "left",
				VAlign:     "top",
				Smart:      true,
				FitIn:      true,
				Filters:    []params.Filter{{Name: "some_filter"}},
			},
		},
	}
	for _, test := range tests {
		if test.name == "" {
			test.name = test.uri
		}
		t.Run(strings.TrimPrefix(test.name, "/"), func(t *testing.T) {
			resp := Parse(test.uri)
			respJSON, _ := json.MarshalIndent(resp, "", "  ")
			expectedJSON, _ := json.MarshalIndent(test.params, "", "  ")
			if !reflect.DeepEqual(resp, test.params) {
				t.Errorf(" = %s, want %s", string(respJSON), string(expectedJSON))
			}
			if test.signer != nil && test.signer.Sign(resp.Path) != resp.Hash {
				t.Errorf("signature mismatch = %s, want %s", resp.Hash, test.signer.Sign(resp.Path))
			}
			if test.params.Hash != "" && test.signer != nil {
				if uri := Generate(test.params, test.signer); uri != test.uri {
					t.Errorf(" = %s, want = %s", uri, test.uri)
				}
			} else if test.params.Unsafe {
				if uri := GenerateUnsafe(test.params); uri != test.uri {
					t.Errorf(" = %s, want = %s", uri, test.uri)
				}
			} else {
				if uri := GeneratePath(test.params); uri != test.uri {
					t.Errorf(" = %s, want = %s", uri, test.uri)
				}
			}
		})
	}
}

func TestParseParamsNegativeDimensionFlip(t *testing.T) {
	assert.Equal(t, "unsafe/-167x-169/foobar", GenerateUnsafe(params.Params{
		Width:  -167,
		Height: -169,
		Image:  "foobar",
	}))
}

func TestParseFilters(t *testing.T) {
	filters, img := parseFilters("filters:watermark(s.glbimg.com/filters:label(abc):watermark(aaa.com/fit-in/filters:aaa(bbb))/aaa.jpg,0,0,0):brightness(-50):grayscale()/some/example/img")
	assert.Equal(t, []params.Filter{
		{"watermark", "s.glbimg.com/filters:label(abc):watermark(aaa.com/fit-in/filters:aaa(bbb))/aaa.jpg,0,0,0"},
		{"brightness", "-50"},
		{"grayscale", ""},
	}, filters)
	assert.Equal(t, "some/example/img", img)

	filters, img = parseFilters("filters:watermark(s.glbimg.com/filters:label(abc):watermark(aaa.com/fit-in/filters:aaa(bbb))/aaa.jpg,0,0,0):brightness(-50):grayscale()")
	assert.Equal(t, []params.Filter{
		{"watermark", "s.glbimg.com/filters:label(abc):watermark(aaa.com/fit-in/filters:aaa(bbb))/aaa.jpg,0,0,0"},
		{"brightness", "-50"},
		{"grayscale", ""},
	}, filters)
	assert.Empty(t, img)

	filters, img = parseFilters("filters:watermark(s.glbimg.com/filters:label(abc):watermark(aaa.com/fit-in/filters:aaa(bbb))/aaa.jpg,0,0,0):brightness(-50):grayscale()/")
	assert.Equal(t, []params.Filter{
		{"watermark", "s.glbimg.com/filters:label(abc):watermark(aaa.com/fit-in/filters:aaa(bbb))/aaa.jpg,0,0,0"},
		{"brightness", "-50"},
		{"grayscale", ""},
	}, filters)
	assert.Empty(t, img)

	filters, img = parseFilters("some/example/img")
	assert.Empty(t, filters)
	assert.Equal(t, "some/example/img", img)

	filters, img = parseFilters("filters:watermark(s.glbimg.com/filters:label(abc):watermark(aaa.com/fit-in/filters:aaa(bbb))/aaa.jpg,0,0,0):format()jpg:brightness(-50):grayscale()")
	assert.Equal(t, []params.Filter{
		{"watermark", "s.glbimg.com/filters:label(abc):watermark(aaa.com/fit-in/filters:aaa(bbb))/aaa.jpg,0,0,0"},
		{"format", ""},
		{"brightness", "-50"},
		{"grayscale", ""},
	}, filters)
	assert.Empty(t, img)
}
