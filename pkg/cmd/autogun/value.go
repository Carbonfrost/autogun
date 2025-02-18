package autogun

import (
	"flag"
	"fmt"
	"strings"

	"github.com/Carbonfrost/autogun/pkg/automation"
)

type Engine int

const (
	Chromedp Engine = iota
)

func (e *Engine) Value() automation.SupportedBinder {
	switch *e {
	case Chromedp:
		return automation.UsingChromedp
	}
	return 0
}

func (*Engine) Synopsis() string {
	return "{chromedp}"
}

func (e *Engine) Set(arg string) error {
	switch strings.ToLower(arg) {
	case "chromedp":
		*e = Chromedp
	default:
		return fmt.Errorf("invalid value: %q", arg)
	}
	return nil
}

func (e Engine) String() string {
	switch e {
	case Chromedp:
		return "chromedp"
	default:
	}
	return ""
}

var _ flag.Value = (*Engine)(nil)
