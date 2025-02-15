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

type selectorTask interface {
	Task
	setSelector(string)
	addSelector(*Selector) error
	setOptions(*Options) error
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
)

func decodeNavigateBlock(block *hcl.Block) (*Navigate, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	f := &Navigate{
		Name:      tryLabel(block, 0),
		DeclRange: block.DefRange,
		NameRange: tryLabelRange(block, 0),
	}

	content, _, moreDiags := block.Body.PartialContent(navigateBlockSchema)
	diags = append(diags, moreDiags...)

	if attr, ok := content.Attributes["url"]; ok {
		f.URL = attr.Expr
	}

	return f, diags
}

func decodeEvalBlock(block *hcl.Block) (*Eval, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	f := &Eval{
		Name:      tryLabel(block, 0),
		DeclRange: block.DefRange,
		NameRange: tryLabelRange(block, 0),
	}

	content, _, moreDiags := block.Body.PartialContent(evalBlockSchema)
	diags = append(diags, moreDiags...)

	if attr, ok := content.Attributes["script"]; ok {
		moreDiags := gohcl.DecodeExpression(attr.Expr, nil, &f.Script)
		diags = append(diags, moreDiags...)
	}

	return f, diags
}

func decodeClickBlock(block *hcl.Block) (*Click, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	f := &Click{
		DeclRange: block.DefRange,
	}

	content, _, moreDiags := block.Body.PartialContent(clickBlockSchema)
	diags = append(diags, moreDiags...)

	moreDiags = andSelectorQueryAttributes(content, f)
	diags = append(diags, moreDiags...)

	return f, diags
}

func decodeWaitVisibleBlock(block *hcl.Block) (*WaitVisible, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	f := &WaitVisible{
		DeclRange: block.DefRange,
	}

	content, _, moreDiags := block.Body.PartialContent(waitVisibleBlockSchema)
	diags = append(diags, moreDiags...)

	moreDiags = andSelectorQueryAttributes(content, f)
	diags = append(diags, moreDiags...)

	return f, diags
}

func decodeScreenshotBlock(block *hcl.Block) (*Screenshot, hcl.Diagnostics) {
	s := &Screenshot{}
	return reduceTask(
		s,
		block,
		supportsDeclRange(&s.DeclRange),
		supportsOptionalLabel(&s.Name, &s.NameRange),
		supportsPartialContentSchema(
			screenshotBlockSchema,
			supportsSelector(s),
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

func (*Navigate) taskSigil()    {}
func (*Eval) taskSigil()        {}
func (*Click) taskSigil()       {}
func (*WaitVisible) taskSigil() {}
func (*Screenshot) taskSigil()  {}

func diagInvalidValue(value string, ty string, subject *hcl.Range) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  "Invalid value",
		Detail:   fmt.Sprintf("The value %q is not supported for %s", value, ty),
		Subject:  subject,
	}
}

var _ Task = (*Navigate)(nil)
var _ Task = (*Eval)(nil)
var _ selectorTask = (*WaitVisible)(nil)
var _ selectorTask = (*Click)(nil)
var _ selectorTask = (*Screenshot)(nil)
