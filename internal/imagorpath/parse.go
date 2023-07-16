package imagorpath

import (
	"fmt"
	"regexp"
	"strconv"
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

var filterRegex = regexp.MustCompile(`^(?P<name>.*)\((?P<args>.+)\)$`)

var regex = regexp.MustCompile(
	`^/?` +
		// params
		`(?:(?P<params>params)/)?` +
		// hmac
		`(?:(?P<hmac>unsafe|[A-Za-z0-9-_=]{8,})/)?` +
		// meta
		`(?:(?P<meta>meta)/)?` +
		// trim
		`(?:(?P<trim>trim)(?:/(?P<trim_by>top-left|bottom-right)(?::(?P<trim_tolerance>\d+))?)?/)?` +
		// crop
		`(?:(?P<crop_left>\d+(?:\.\d+)?)x(?P<crop_top>\d+(?:\.\d+)?)(?::(?P<crop_right>\d+(?:\.\d+)?)x(?P<crop_bottom>\d+(?:\.\d+)?))?/)?` +
		// fit-in
		`(?:(?P<fit_in>fit-in)/)?` +
		// stretch
		`(?:(?P<stretch>stretch)/)?` +
		// dimensions
		`(?:(?P<h_flip>-)?(?P<width>\d+)x(?P<v_flip>-)?(?P<height>\d+)/)?` +
		// paddings
		`(?:(?P<padding_left>\d+)x(?P<padding_top>\d+)(?::(?P<padding_right>\d+)x(?P<padding_bottom>\d+))?/)?` +
		// h_align
		`(?:(?P<h_align>left|right|center)/)?` +
		// v_align
		`(?:(?P<v_align>top|bottom|middle)/)?` +
		// smart
		`(?:(?P<smart>smart)/)?` +
		// filters
		`(?:filters:(?P<filters>.*)/)?` +
		// image
		`(?P<image>.+)?` +
		`$`,
)

type parser struct {
	params *params.Params
}

func (p parser) parseParams(value string) error {
	p.params.Echo = value == "params"
	return nil
}

func (p parser) parseHMAC(value string) error {
	if value == "unsafe" {
		p.params.Unsafe = true
	} else if value != "" {
		p.params.Hash = value
	}
	return nil
}

func (p parser) parseMeta(value string) error {
	p.params.Meta = value == "meta"
	return nil
}

func (p parser) parseTrim(value string) error {
	p.params.Trim = value == "trim"
	return nil
}

func (p parser) parseTrimBy(value string) error {
	if value == params.TrimByTopLeft || value == params.TrimByBottomRight {
		p.params.TrimBy = value
	}
	return nil
}

func (p parser) parseTrimTolerance(value string) error {
	t, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid trim_tolerance value: %s", value)
	}
	p.params.TrimTolerance = t
	return nil
}

func (p parser) parseCropLeft(value string) error {
	c, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("invalid crop_left value: %s", value)
	}
	p.params.CropLeft = c
	return nil
}

func (p parser) parseCropTop(value string) error {
	c, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("invalid crop_top value: %s", value)
	}
	p.params.CropTop = c
	return nil
}

func (p parser) parseCropRight(value string) error {
	c, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("invalid crop_right value: %s", value)
	}
	p.params.CropRight = c
	return nil
}

func (p parser) parseCropBottom(value string) error {
	c, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("invalid crop_bottom value: %s", value)
	}
	p.params.CropBottom = c
	return nil
}

func (p parser) parseFitIn(value string) error {
	p.params.FitIn = value == "fit-in"
	return nil
}

func (p parser) parseStretch(value string) error {
	p.params.Stretch = value == "stretch"
	return nil
}

func (p parser) parseHFlip(value string) error {
	p.params.HFlip = value == "-"
	return nil
}

func (p parser) parseWidth(value string) error {
	w, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid width value: %s", value)
	}
	p.params.Width = w
	return nil
}

func (p parser) parseVFlip(value string) error {
	p.params.VFlip = value == "-"
	return nil
}

func (p parser) parseHeight(value string) error {
	h, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid height value: %s", value)
	}
	p.params.Height = h
	return nil
}

func (p parser) parsePaddingLeft(value string) error {
	d, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid padding_left value: %s", value)
	}
	p.params.PaddingLeft = d
	return nil
}

func (p parser) parsePaddingTop(value string) error {
	d, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid padding_top value: %s", value)
	}
	p.params.PaddingTop = d
	return nil
}

func (p parser) parsePaddingRight(value string) error {
	d, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid padding_right value: %s", value)
	}
	p.params.PaddingRight = d
	return nil
}

func (p parser) parsePaddingBottom(value string) error {
	d, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid padding_bottom value: %s", value)
	}
	p.params.PaddingBottom = d
	return nil
}

func (p parser) parseHAlign(value string) error {
	if value != params.HAlignLeft && value != params.HAlignRight {
		return fmt.Errorf("invalid h_align value: %s", value)
	}
	p.params.HAlign = value
	return nil
}

func (p parser) parseVAlign(value string) error {
	if value != params.VAlignTop && value != params.VAlignBottom {
		return fmt.Errorf("invalid v_align value: %s", value)
	}
	p.params.VAlign = value
	return nil
}

func (p parser) parseSmart(value string) error {
	p.params.Smart = value == "smart"
	return nil
}

func (p parser) parseFilters(value string) error {
	rawFilters := strings.Split(value, ":")
	filters := make([]params.Filter, len(rawFilters))
	for i, rawFilter := range rawFilters {
		filter := params.Filter{}
		matches := filterRegex.FindStringSubmatch(rawFilter)
		for i, group := range filterRegex.SubexpNames() {
			switch group {
			case "name":
				filter.Name = matches[i]
			case "args":
				filter.Args = matches[i]
			}
		}
		filters[i] = filter
	}
	p.params.Filters = filters
	return nil
}

func (p parser) parseImage(image string) error {
	p.params.Image = image
	return nil
}

func (p parser) parse(path string) error {
	path = breaksCleaner.Replace(path)

	// TODO: add timeout
	matches := regex.FindStringSubmatch(path)
	if len(matches) <= 1 {
		return fmt.Errorf("invalid path: %s", path)
	}

	for i, group := range regex.SubexpNames() {
		fmt.Printf("%s: %s\n", group, matches[i])
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
		case "trim_by":
			if err := p.parseTrimBy(matches[i]); err != nil {
				return err
			}
		case "trim_tolerance":
			if err := p.parseTrimTolerance(matches[i]); err != nil {
				return err
			}
		case "crop_left":
			if err := p.parseCropLeft(matches[i]); err != nil {
				return err
			}
		case "crop_top":
			if err := p.parseCropTop(matches[i]); err != nil {
				return err
			}
		case "crop_right":
			if err := p.parseCropRight(matches[i]); err != nil {
				return err
			}
		case "crop_bottom":
			if err := p.parseCropBottom(matches[i]); err != nil {
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
		case "h_flip":
			if err := p.parseHFlip(matches[i]); err != nil {
				return err
			}
		case "width":
			if err := p.parseWidth(matches[i]); err != nil {
				return err
			}
		case "v_flip":
			if err := p.parseVFlip(matches[i]); err != nil {
				return err
			}
		case "height":
			if err := p.parseHeight(matches[i]); err != nil {
				return err
			}
		case "padding_left":
			if err := p.parsePaddingLeft(matches[i]); err != nil {
				return err
			}
		case "padding_top":
			if err := p.parsePaddingTop(matches[i]); err != nil {
				return err
			}
		case "padding_right":
			if err := p.parsePaddingRight(matches[i]); err != nil {
				return err
			}
		case "padding_bottom":
			if err := p.parsePaddingBottom(matches[i]); err != nil {
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
