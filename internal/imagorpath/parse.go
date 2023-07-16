package imagorpath

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cshum/imagor/pkg/params"
)

var breaksCleaner = strings.NewReplacer(
	"\r\n", "",
	"\r", "",
	"\n", "",
	"\v", "",
	"\f", "",
	"\u0085", "",
	"\u2028", "",
	"\u2029", "",
)

var regex = regexp.MustCompile(
	`^/?` +
		// params
		`(?P<params>params/)?` +
		// hmac
		`(?P<hmac>(unsafe/)|([A-Za-z0-9-_=]{8,})/)?` +
		// meta
		`(?P<meta>meta/)?` +
		// trim
		`(?P<trim>trim/(:(?P<trim_align>(top-left|bottom-right))?(:(?P<trim_tolerance>\d+))?/)?)?` +
		// crop
		`(?P<crop>((0?\.)?\d+)x((0?\.)?\d+):((?P<crop_x>((0?\.)?\d+))x((?P<crop_y>((0?\.)?\d+))))/)?` +
		// fit-in
		`(?P<fit_in>fit-in/)?` +
		// stretch
		`(?P<stretch>stretch/)?` +
		// dimensions
		`(?P<dimensions>((\-?)(\d*)x(\-?)(\d*)/)?)` +
		// paddings
		`(?P<paddings>((\d+)x(\d+)(:(\d+)x(\d+))?/)?)?` +
		// h_align
		`(?P<h_align>(left|right|center)/)?` +
		// v_align
		`(?P<v_align>(top|bottom|middle)/)?` +
		// smart
		`(?P<smart>smart/)?` +
		// filters
		`(?P<filters>.*/)?` +
		// image
		`(?P<image>.+)?` +
		`$`,
)

type parser struct {
	params *params.Params
}

func (p parser) parseParams(_ string) error {
	p.params.Echo = true
	return nil
}

func (p parser) parseHMAC(hmac string) error {
	p.params.Unsafe = hmac == "unsafe/"
	p.params.Hash = hmac
	return nil
}

func (p parser) parseMeta(_ string) error {
	p.params.Meta = true
	return nil
}

func (p parser) parseTrim(trim string) error {
	// TODO
	return nil
}

func (p parser) parseCrop(crop string) error {
	// TODO
	return nil
}

func (p parser) parseFitIn(_ string) error {
	p.params.FitIn = true
	return nil
}

func (p parser) parseStretch(_ string) error {
	p.params.Stretch = true
	return nil
}

func (p parser) parseDimensions(dimensions string) error {
	// TODO
	return nil
}

func (p parser) parsePaddings(paddings string) error {
	// TODO
	return nil
}

func (p parser) parseHAlign(hAlign string) error {
	// TODO
	return nil
}

func (p parser) parseVAlign(vAlign string) error {
	// TODO
	return nil
}

func (p parser) parseSmart(_ string) error {
	p.params.Smart = true
	return nil
}

func (p parser) parseFilters(filters string) error {
	// TODO
	return nil
}

func (p parser) parseImage(image string) error {
	p.params.Image = image
	return nil
}

func (p parser) parse(path string) error {
	path = breaksCleaner.Replace(path)

	matches := regex.FindStringSubmatch(path)
	if len(matches) <= 1 {
		return fmt.Errorf("invalid path: %s", path)
	}

	for i, group := range regex.SubexpNames() {
		switch group {
		case "params":
			if err := p.parseParams(matches[i]); err != nil {
				return err
			}
		case "hmac":
			if err := p.parseHMAC(matches[i]); err != nil {
				return err
			}
		case "meta":
			if err := p.parseMeta(matches[i]); err != nil {
				return err
			}
		case "trim":
			if err := p.parseTrim(matches[i]); err != nil {
				return err
			}
		case "crop":
			if err := p.parseCrop(matches[i]); err != nil {
				return err
			}
		case "fit_in":
			if err := p.parseFitIn(matches[i]); err != nil {
				return err
			}
		case "stretch":
			if err := p.parseStretch(matches[i]); err != nil {
				return err
			}
		case "dimensions":
			if err := p.parseDimensions(matches[i]); err != nil {
				return err
			}
		case "paddings":
			if err := p.parsePaddings(matches[i]); err != nil {
				return err
			}
		case "h_align":
			if err := p.parseHAlign(matches[i]); err != nil {
				return err
			}
		case "v_align":
			if err := p.parseVAlign(matches[i]); err != nil {
				return err
			}
		case "smart":
			if err := p.parseSmart(matches[i]); err != nil {
				return err
			}
		case "filters":
			if err := p.parseFilters(matches[i]); err != nil {
				return err
			}
		case "image":
			if err := p.parseImage(matches[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

// Parse Params struct from imagorpath
func Parse(path string) (params.Params, error) {
	params := params.Params{}
	parser := parser{params: &params}
	err := parser.parse(path)
	return *parser.params, err
}
