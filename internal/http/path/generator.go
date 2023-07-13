package path

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/cshum/imagor/pkg/params"
)

func escapeImage(image string) string {
	if strings.Contains(image, "?") ||
		strings.HasPrefix(image, "trim/") ||
		strings.HasPrefix(image, "meta/") ||
		strings.HasPrefix(image, "fit-in/") ||
		strings.HasPrefix(image, "stretch/") ||
		strings.HasPrefix(image, "top/") ||
		strings.HasPrefix(image, "left/") ||
		strings.HasPrefix(image, "right/") ||
		strings.HasPrefix(image, "bottom/") ||
		strings.HasPrefix(image, "center/") ||
		strings.HasPrefix(image, "smart/") {
		return url.QueryEscape(image)
	}
	return image
}

type generator struct {
	segments []string
}

func (g *generator) GenerateRaw(p params.Params) string {
	path := g.generate(p)
	return path
}

func (g *generator) GenerateUnsafe(p params.Params) string {
	path := g.GenerateRaw(p)
	return "unsafe/" + path
}

func (g *generator) GenerateSigned(p params.Params, signer Signer) string {
	path := g.GenerateRaw(p)
	return signer.Sign(path) + "/" + path
}

func (g *generator) generate(p params.Params) string {
	if p.ShouldMeta() {
		g.segments = append(g.segments, "meta")
	}
	if p.ShouldTrim() {
		trims := []string{"trim"}
		if p.TrimBy == params.TrimByBottomRight {
			trims = append(trims, "bottom-right")
		}
		if p.TrimTolerance > 0 {
			trims = append(trims, strconv.Itoa(p.TrimTolerance))
		}
		g.segments = append(g.segments, strings.Join(trims, ":"))
	}
	if p.ShouldCrop() {
		g.segments = append(g.segments, fmt.Sprintf(
			"%sx%s:%sx%s",
			strconv.FormatFloat(p.CropLeft, 'f', -1, 64),
			strconv.FormatFloat(p.CropTop, 'f', -1, 64),
			strconv.FormatFloat(p.CropRight, 'f', -1, 64),
			strconv.FormatFloat(p.CropBottom, 'f', -1, 64)))
	}
	if p.ShouldFitIn() {
		g.segments = append(g.segments, "fit-in")
	}
	if p.ShouldStretch() {
		g.segments = append(g.segments, "stretch")
	}
	if p.ShouldFlip() {
		if p.Width < 0 {
			p.HFlip = !p.HFlip
			p.Width = -p.Width
		}
		if p.Height < 0 {
			p.VFlip = !p.VFlip
			p.Height = -p.Height
		}

		var hFlipStr, vFlipStr string
		if p.HFlip {
			hFlipStr = "-"
		}
		if p.VFlip {
			vFlipStr = "-"
		}
		g.segments = append(g.segments, fmt.Sprintf(
			"%s%dx%s%d", hFlipStr, p.Width, vFlipStr, p.Height,
		))
	}
	if p.ShouldPad() {
		if p.PaddingLeft == p.PaddingRight && p.PaddingTop == p.PaddingBottom {
			g.segments = append(g.segments, fmt.Sprintf("%dx%d", p.PaddingLeft, p.PaddingTop))
		} else {
			g.segments = append(g.segments, fmt.Sprintf(
				"%dx%d:%dx%d",
				p.PaddingLeft, p.PaddingTop,
				p.PaddingRight, p.PaddingBottom,
			))
		}
	}
	if p.ShouldHAlign() {
		g.segments = append(g.segments, p.HAlign)
	}
	if p.ShouldVAlign() {
		g.segments = append(g.segments, p.VAlign)
	}
	if p.ShouldSmart() {
		g.segments = append(g.segments, "smart")
	}
	if p.HasFilters() {
		filters := make([]string, len(p.Filters))
		for i, f := range p.Filters {
			filters[i] = fmt.Sprintf("%s(%s)", f.Name, f.Args)
		}
		g.segments = append(g.segments, "filters:"+strings.Join(filters, ":"))
	}
	image := escapeImage(p.Image)
	g.segments = append(g.segments, image)

	return strings.Join(g.segments, "/")
}

type Signer interface {
	Sign(string) string
}

func New() *generator {
	return &generator{}
}
