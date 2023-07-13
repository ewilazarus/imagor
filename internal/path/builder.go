package path

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/cshum/imagor/pkg/params"
)

type builder struct {
	params.Params
	segments []string
}

func (b builder) appendMeta() builder {
	if b.Meta {
		b.segments = append(b.segments, "meta")
	}
	return b
}

func (b builder) appendTrim() builder {
	if b.Trim || (b.TrimBy == params.TrimByTopLeft || b.TrimBy == params.TrimByBottomRight) {
		trims := []string{"trim"}
		if b.TrimBy == params.TrimByBottomRight {
			trims = append(trims, "bottom-right")
		}
		if b.TrimTolerance > 0 {
			trims = append(trims, strconv.Itoa(b.TrimTolerance))
		}
		b.segments = append(b.segments, strings.Join(trims, ":"))
	}
	return b
}

func (b builder) appendCrop() builder {
	if b.CropTop > 0 || b.CropRight > 0 || b.CropLeft > 0 || b.CropBottom > 0 {
		b.segments = append(b.segments, fmt.Sprintf(
			"%sx%s:%sx%s",
			strconv.FormatFloat(b.CropLeft, 'f', -1, 64),
			strconv.FormatFloat(b.CropTop, 'f', -1, 64),
			strconv.FormatFloat(b.CropRight, 'f', -1, 64),
			strconv.FormatFloat(b.CropBottom, 'f', -1, 64)))
	}
	return b
}

func (b builder) appendFitIn() builder {
	if b.FitIn {
		b.segments = append(b.segments, "fit-in")
	}
	return b
}

func (b builder) appendStretch() builder {
	if b.Stretch {
		b.segments = append(b.segments, "stretch")
	}
	return b
}

func (b builder) appendDimensions() builder {
	if b.HFlip || b.Width != 0 || b.VFlip || b.Height != 0 || b.PaddingLeft > 0 || b.PaddingTop > 0 {
		if b.Width < 0 {
			b.HFlip = !b.HFlip
			b.Width = -b.Width
		}
		if b.Height < 0 {
			b.VFlip = !b.VFlip
			b.Height = -b.Height
		}

		var hFlipStr, vFlipStr string
		if b.HFlip {
			hFlipStr = "-"
		}
		if b.VFlip {
			vFlipStr = "-"
		}

		dimensions := fmt.Sprintf("%s%dx%s%d", hFlipStr, b.Width, vFlipStr, b.Height)
		b.segments = append(b.segments, dimensions)
	}
	return b
}

func (b builder) appendPaddings() builder {
	if b.PaddingLeft > 0 || b.PaddingTop > 0 || b.PaddingRight > 0 || b.PaddingBottom > 0 {
		if b.PaddingLeft == b.PaddingRight && b.PaddingTop == b.PaddingBottom {
			b.segments = append(b.segments, fmt.Sprintf("%dx%d", b.PaddingLeft, b.PaddingTop))
		} else {
			b.segments = append(b.segments, fmt.Sprintf(
				"%dx%d:%dx%d",
				b.PaddingLeft, b.PaddingTop,
				b.PaddingRight, b.PaddingBottom,
			))
		}
	}
	return b
}

func (b builder) appendHAlign() builder {
	if b.HAlign == params.HAlignLeft || b.HAlign == params.HAlignRight {
		b.segments = append(b.segments, b.HAlign)
	}
	return b
}

func (b builder) appendVAlign() builder {
	if b.VAlign == params.VAlignTop || b.VAlign == params.VAlignBottom {
		b.segments = append(b.segments, b.VAlign)
	}
	return b
}

func (b builder) appendSmart() builder {
	if b.Smart {
		b.segments = append(b.segments, "smart")
	}
	return b
}

func (b builder) appendFilters() builder {
	filters := make([]string, len(b.Filters))
	for i, f := range b.Filters {
		filters[i] = fmt.Sprintf("%s(%s)", f.Name, f.Args)
	}
	b.segments = append(b.segments, "filters:"+strings.Join(filters, ":"))
	return b
}

func (b builder) appendImage() builder {
	image := b.Image
	if strings.Contains(b.Image, "?") ||
		strings.HasPrefix(b.Image, "trim/") ||
		strings.HasPrefix(b.Image, "meta/") ||
		strings.HasPrefix(b.Image, "fit-in/") ||
		strings.HasPrefix(b.Image, "stretch/") ||
		strings.HasPrefix(b.Image, "top/") ||
		strings.HasPrefix(b.Image, "left/") ||
		strings.HasPrefix(b.Image, "right/") ||
		strings.HasPrefix(b.Image, "bottom/") ||
		strings.HasPrefix(b.Image, "center/") ||
		strings.HasPrefix(b.Image, "smart/") {
		image = url.QueryEscape(b.Image)
	}
	b.segments = append(b.segments, image)
	return b
}

func (b builder) stringify() string {
	return strings.Join(b.segments, "/")
}

func (b builder) BuildUnprefixed() string {
	// Note that the order matters
	return b.
		appendMeta().
		appendTrim().
		appendCrop().
		appendFitIn().
		appendStretch().
		appendDimensions().
		appendPaddings().
		appendHAlign().
		appendVAlign().
		appendSmart().
		appendFilters().
		appendImage().
		stringify()
}

func (b builder) BuildUnsafe() string {
	return "unsafe/" + b.BuildUnprefixed()
}

type Signer interface {
	Sign(string) string
}

func (b builder) BuildSigned(signer Signer) string {
	path := b.BuildUnprefixed()
	return signer.Sign(path) + "/" + path
}

func NewBuilder(params params.Params) builder {
	return builder{Params: params}
}
