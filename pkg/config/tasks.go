package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/hcl/v2"
)

type Task interface {
	taskSigil()
}

type Navigate struct {
	DeclRange hcl.Range
	NameRange hcl.Range
	Name      string
	URL       hcl.Expression
}

type NavigateForward struct {
	DeclRange hcl.Range
}

type NavigateBack struct {
	DeclRange hcl.Range
}

type Eval struct {
	DeclRange hcl.Range
	NameRange hcl.Range
	Name      string
	Script    string
}

type Click struct {
	DeclRange hcl.Range
	Selector  string
	Selectors []*Selector
	Options   *Options
}

type WaitVisible struct {
	DeclRange hcl.Range
	Selector  string
	Selectors []*Selector
	Options   *Options
}

type Screenshot struct {
	DeclRange hcl.Range
	NameRange hcl.Range
	Name      string
	Selector  string
	Selectors []*Selector
	Options   *Options
}

type Options struct {
	DeclRange     hcl.Range
	RetryInterval *time.Duration
	AtLeast       *int
}

type Sleep struct {
	DeclRange hcl.Range
	Duration  time.Duration
}

type Reload struct {
	DeclRange hcl.Range
}

type Stop struct {
	DeclRange hcl.Range
}

var (
	navigateBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "url"},
		},
	}

	evalBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "script"},
		},
	}

	clickBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "selector"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "selector"},
			{Type: "options"},
		},
	}

	waitVisibleBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "selector"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "selector"},
			{Type: "options"},
		},
	}

	screenshotBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "selector"},
		},
		Blocks: []hcl.BlockHeaderSchema{
			{Type: "selector"},
			{Type: "options"},
		},
	}

	optionsBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "at_least"},
			{Name: "retry_interval"},
		},
	}

	sleepBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "duration"},
		},
	}

	navigateForwardBlockSchema = &hcl.BodySchema{}
	navigateBackBlockSchema    = &hcl.BodySchema{}
	reloadBlockSchema          = &hcl.BodySchema{}
	stopBlockSchema            = &hcl.BodySchema{}
)

func decodeNavigateBlock(block *hcl.Block) (*Navigate, hcl.Diagnostics) {
	f := new(Navigate)
	return reduceTask(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsOptionalLabel(&f.Name, &f.NameRange),
		supportsPartialContentSchema(
			navigateBlockSchema,
			withAttributeExpression("url", &f.URL),
		),
	)
}

func decodeNavigateForwardBlock(block *hcl.Block) (*NavigateForward, hcl.Diagnostics) {
	f := new(NavigateForward)
	return reduceTask(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsPartialContentSchema(
			navigateForwardBlockSchema,
		),
	)
}

func decodeNavigateBackBlock(block *hcl.Block) (*NavigateBack, hcl.Diagnostics) {
	f := new(NavigateBack)
	return reduceTask(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsPartialContentSchema(
			navigateBackBlockSchema,
		),
	)
}

func decodeReloadBlock(block *hcl.Block) (*Reload, hcl.Diagnostics) {
	f := new(Reload)
	return reduceTask(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsPartialContentSchema(
			reloadBlockSchema,
		),
	)
}

func decodeStopBlock(block *hcl.Block) (*Stop, hcl.Diagnostics) {
	f := new(Stop)
	return reduceTask(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsPartialContentSchema(
			stopBlockSchema,
		),
	)
}

func decodeEvalBlock(block *hcl.Block) (*Eval, hcl.Diagnostics) {
	f := new(Eval)
	return reduceTask(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsOptionalLabel(&f.Name, &f.NameRange),
		supportsPartialContentSchema(
			evalBlockSchema,
			withAttribute("script", &f.Script),
		),
	)
}

func decodeClickBlock(block *hcl.Block) (*Click, hcl.Diagnostics) {
	f := new(Click)
	return reduceTask(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsPartialContentSchema(
			clickBlockSchema,
			withAttribute("selector", &f.Selector),
			supportsSelectorBlocks(&f.Selectors, &f.Options),
		),
	)
}

func decodeSleepBlock(block *hcl.Block) (*Sleep, hcl.Diagnostics) {
	f := new(Sleep)
	return reduceTask(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsPartialContentSchema(
			sleepBlockSchema,
			withAttributeParser("duration", f.setDuration, time.ParseDuration),
		),
	)
}

func decodeWaitVisibleBlock(block *hcl.Block) (*WaitVisible, hcl.Diagnostics) {
	f := new(WaitVisible)
	return reduceTask(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsPartialContentSchema(
			waitVisibleBlockSchema,
			withAttribute("selector", &f.Selector),
			supportsSelectorBlocks(&f.Selectors, &f.Options),
		),
	)
}

func decodeScreenshotBlock(block *hcl.Block) (*Screenshot, hcl.Diagnostics) {
	s := new(Screenshot)
	return reduceTask(
		s,
		block,
		supportsDeclRange(&s.DeclRange),
		supportsOptionalLabel(&s.Name, &s.NameRange),
		supportsPartialContentSchema(
			screenshotBlockSchema,
			withAttribute("selector", &s.Selector),
			supportsSelectorBlocks(&s.Selectors, &s.Options),
		),
	)
}

func decodeOptionsBlock(block *hcl.Block) (*Options, hcl.Diagnostics) {
	s := new(Options)
	return reduce(
		s,
		block,
		supportsDeclRange(&s.DeclRange),
		supportsPartialContentSchema(
			optionsBlockSchema,
			withAttributeParser("retry_interval", s.setRetryInterval, time.ParseDuration),
			withAttributeParser("at_least", s.setAtLeast, strconv.Atoi),
		),
	)
}

func (o *Options) setRetryInterval(n time.Duration) {
	o.RetryInterval = &n
}

func (o *Options) setAtLeast(n int) {
	o.AtLeast = &n
}

func (o *Sleep) setDuration(n time.Duration) {
	o.Duration = n
}

func (*Automation) taskSigil()      {}
func (*Click) taskSigil()           {}
func (*Eval) taskSigil()            {}
func (*Navigate) taskSigil()        {}
func (*NavigateBack) taskSigil()    {}
func (*NavigateForward) taskSigil() {}
func (*Reload) taskSigil()          {}
func (*Screenshot) taskSigil()      {}
func (*Sleep) taskSigil()           {}
func (*Stop) taskSigil()            {}
func (*WaitVisible) taskSigil()     {}

func diagInvalidValue(value string, ty string, subject *hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Invalid value",
		Detail:   fmt.Sprintf("The value %q is not supported for %s", value, ty),
		Subject:  subject,
	}
}
