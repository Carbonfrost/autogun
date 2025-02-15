package config

import (
	"fmt"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
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

var (
	navigateBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "url"},
		},
		Blocks: []hcl.BlockHeaderSchema{},
	}

	evalBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "script"},
		},
		Blocks: []hcl.BlockHeaderSchema{},
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
		Blocks: []hcl.BlockHeaderSchema{},
	}

	navigateForwardBlockSchema = &hcl.BodySchema{}
	navigateBackBlockSchema    = &hcl.BodySchema{}
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

func decodeWaitVisibleBlock(block *hcl.Block) (*WaitVisible, hcl.Diagnostics) {
	f := new(WaitVisible)
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
	var diags hcl.Diagnostics
	f := &Options{
		DeclRange: block.DefRange,
	}

	content, _, moreDiags := block.Body.PartialContent(optionsBlockSchema)
	diags = append(diags, moreDiags...)
	if attr, ok := content.Attributes["retry_interval"]; ok {
		var text string
		moreDiags := gohcl.DecodeExpression(attr.Expr, nil, &text)
		if text != "" {
			dur, err := time.ParseDuration(text)
			f.RetryInterval = &dur
			if err != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "Cannot convert duration",
					Detail:   fmt.Sprintf("Cannot convert duration: %s", err.Error()),
					Subject:  attr.Expr.StartRange().Ptr(),
					Context:  attr.Expr.Range().Ptr(),
				})
			}
		}
		diags = append(diags, moreDiags...)
	}
	if attr, ok := content.Attributes["at_least"]; ok {
		moreDiags := gohcl.DecodeExpression(attr.Expr, nil, &f.AtLeast)
		diags = append(diags, moreDiags...)
	}

	return f, diags
}

func (c *Click) setSelector(s string) {
	c.Selector = s
}

func (c *Click) addSelector(s *Selector) error {
	c.Selectors = append(c.Selectors, s)
	return nil
}

func (c *Click) setOptions(o *Options) error {
	c.Options = o
	return nil
}

func (w *WaitVisible) setSelector(s string) {
	w.Selector = s
}

func (w *WaitVisible) addSelector(s *Selector) error {
	w.Selectors = append(w.Selectors, s)
	return nil
}

func (w *WaitVisible) setOptions(o *Options) error {
	w.Options = o
	return nil
}

func (s *Screenshot) setSelector(t string) {
	s.Selector = t
}

func (s *Screenshot) addSelector(t *Selector) error {
	s.Selectors = append(s.Selectors, t)
	return nil
}

func (s *Screenshot) setOptions(o *Options) error {
	s.Options = o
	return nil
}

func (*Automation) taskSigil()      {}
func (*Navigate) taskSigil()        {}
func (*NavigateForward) taskSigil() {}
func (*NavigateBack) taskSigil()    {}
func (*Eval) taskSigil()            {}
func (*Click) taskSigil()           {}
func (*WaitVisible) taskSigil()     {}
func (*Screenshot) taskSigil()      {}

func diagInvalidValue(value string, ty string, subject *hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Invalid value",
		Detail:   fmt.Sprintf("The value %q is not supported for %s", value, ty),
		Subject:  subject,
	}
}

var (
	_ Task = (*Navigate)(nil)
	_ Task = (*NavigateForward)(nil)
	_ Task = (*NavigateBack)(nil)
	_ Task = (*Eval)(nil)
	_ Task = (*Automation)(nil)
	_ Task = (*WaitVisible)(nil)
	_ Task = (*Click)(nil)
	_ Task = (*Screenshot)(nil)
)
