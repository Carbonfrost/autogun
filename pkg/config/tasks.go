package config

import (
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
	URL       string
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
}

type WaitVisible struct {
	DeclRange hcl.Range
	Selector  string
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
		Blocks: []hcl.BlockHeaderSchema{},
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
		moreDiags := gohcl.DecodeExpression(attr.Expr, nil, &f.URL)
		diags = append(diags, moreDiags...)
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

	moreDiags = andSelectorQueryAttributes(content, &f.Selector)
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

	moreDiags = andSelectorQueryAttributes(content, &f.Selector)
	diags = append(diags, moreDiags...)

	return f, diags
}

func andSelectorQueryAttributes(content *hcl.BodyContent, dstSelector *string) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if attr, ok := content.Attributes["selector"]; ok {
		moreDiags := gohcl.DecodeExpression(attr.Expr, nil, dstSelector)
		diags = append(diags, moreDiags...)
	}

	return diags
}

func (*Navigate) taskSigil()    {}
func (*Eval) taskSigil()        {}
func (*Click) taskSigil()       {}
func (*WaitVisible) taskSigil() {}

var _ Task = (*Navigate)(nil)
var _ Task = (*Eval)(nil)
