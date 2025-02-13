package config

import (
	"fmt"

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
}

type WaitVisible struct {
	DeclRange hcl.Range
	Selector  string
	Selectors []*Selector
}

type selectorAction interface {
	Task
	setSelector(string)
	addSelector(*Selector) error
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
		},
	}

	waitVisibleBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "selector"},
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

func (c *Click) setSelector(s string) {
	c.Selector = s
}

func (c *Click) addSelector(s *Selector) error {
	c.Selectors = append(c.Selectors, s)
	return nil
}

func (w *WaitVisible) setSelector(s string) {
	w.Selector = s
}

func (w *WaitVisible) addSelector(s *Selector) error {
	w.Selectors = append(w.Selectors, s)
	return nil
}

func (*Navigate) taskSigil()    {}
func (*Eval) taskSigil()        {}
func (*Click) taskSigil()       {}
func (*WaitVisible) taskSigil() {}

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
var _ selectorAction = (*WaitVisible)(nil)
var _ selectorAction = (*Click)(nil)
