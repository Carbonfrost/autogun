package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

type mapper func(*hcl.Block) hcl.Diagnostics
type partialContentMapper func(*hcl.BodyContent) hcl.Diagnostics

func reduceTask[T Task](a T, block *hcl.Block, mappers ...mapper) (T, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	for _, m := range mappers {
		diags = append(diags, m(block)...)
	}
	return a, diags
}

func supportsDeclRange(d *hcl.Range) mapper {
	return func(block *hcl.Block) hcl.Diagnostics {
		*d = block.DefRange
		return nil
	}
}

func supportsOptionalLabel(name *string, nameRange *hcl.Range) mapper {
	return func(block *hcl.Block) hcl.Diagnostics {
		*name = tryLabel(block, 0)
		*nameRange = block.DefRange
		return nil
	}
}

func supportsPartialContentSchema(schema *hcl.BodySchema, att ...partialContentMapper) mapper {
	return func(block *hcl.Block) hcl.Diagnostics {
		content, _, diags := block.Body.PartialContent(schema)
		for _, a := range att {
			diags = append(diags, a(content)...)
		}
		return diags
	}
}

func supportsSelector(dst selectorTask) partialContentMapper {
	return func(content *hcl.BodyContent) hcl.Diagnostics {
		return andSelectorQueryAttributes(content, dst)
	}
}

func withAttribute(name string, value any) partialContentMapper {
	return func(content *hcl.BodyContent) hcl.Diagnostics {
		if attr, ok := content.Attributes[name]; ok {
			return gohcl.DecodeExpression(attr.Expr, nil, value)
		}
		return nil
	}
}
