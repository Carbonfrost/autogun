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

func (*Navigate) taskSigil() {}
func (*Eval) taskSigil()     {}

var _ Task = (*Navigate)(nil)
var _ Task = (*Eval)(nil)
