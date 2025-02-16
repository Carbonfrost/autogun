package config

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

type mapper func(*hcl.Block) hcl.Diagnostics
type partialContentMapper func(*hcl.BodyContent) hcl.Diagnostics
type blockMapping[T any] map[string]func(*hcl.Block) (T, hcl.Diagnostics)

func reduce[T any](a T, block *hcl.Block, mappers ...mapper) (T, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	for _, m := range mappers {
		diags = append(diags, m(block)...)
	}
	return a, diags
}

func reduceTask[T Task](a T, block *hcl.Block, mappers ...mapper) (T, hcl.Diagnostics) {
	return reduce(a, block, mappers...)
}

func supportsDeclRange(d *hcl.Range) mapper {
	return func(block *hcl.Block) hcl.Diagnostics {
		*d = block.DefRange
		return nil
	}
}

func supportsOptionalLabel(name *string, nameRange *hcl.Range) mapper {
	return func(block *hcl.Block) hcl.Diagnostics {
		label := tryLabel(block, 0)
		*name = label
		if !validIdentifier(label) {
			return hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Invalid identifier name",
					Detail:   badIdentifierDetail,
				},
			}
		}
		*nameRange = tryLabelRange(block, 0)
		return nil
	}
}

func appendsTo[T any, Slice ~[]T](target *Slice, m blockMapping[T]) partialContentMapper {
	return func(content *hcl.BodyContent) hcl.Diagnostics {
		var diags hcl.Diagnostics
		var results []T
		for _, block := range content.Blocks {
			cfg, cfgDiags := m[block.Type](block)
			diags = append(diags, cfgDiags...)
			if any(cfg) != nil {
				results = append(results, cfg)
			}
		}
		*target = append(*target, results...)
		return diags
	}
}

// contravariant conversion of return type
func taskMapping[T Task](fn func(*hcl.Block) (T, hcl.Diagnostics)) func(*hcl.Block) (Task, hcl.Diagnostics) {
	return func(b *hcl.Block) (Task, hcl.Diagnostics) {
		return fn(b)
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

func withAttr(name string, fn func(*hcl.Attribute) hcl.Diagnostics) partialContentMapper {
	return func(content *hcl.BodyContent) hcl.Diagnostics {
		if attr, ok := content.Attributes[name]; ok {
			return fn(attr)
		}
		return nil
	}
}

func withAttribute(name string, value any) partialContentMapper {
	return withAttr(name, func(attr *hcl.Attribute) hcl.Diagnostics {
		return gohcl.DecodeExpression(attr.Expr, nil, value)
	})
}

func withAttributeExpression(name string, value *hcl.Expression) partialContentMapper {
	return withAttr(name, func(attr *hcl.Attribute) hcl.Diagnostics {
		*value = attr.Expr
		return nil
	})
}

func withAttributeParser[T any](name string, valueThunk func(T), parser func(string) (T, error)) partialContentMapper {
	// A valueThunk is used instead of setting the variable directly because
	// this function can _also_ be used for attributes that are optional
	// and represented in their models using pointers (instead of values directly).
	// An example of this is Options.RetryInterval which is *time.Duration
	return withAttr(name, func(attr *hcl.Attribute) hcl.Diagnostics {
		var text string

		diags := gohcl.DecodeExpression(attr.Expr, nil, &text)
		dur, err := parser(text)
		if err == nil {
			valueThunk(dur)
		} else {
			var it T
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Cannot convert %T", it),
				Detail:   fmt.Sprintf("Cannot convert %T: %s", it, err.Error()),
				Subject:  attr.Expr.StartRange().Ptr(),
				Context:  attr.Expr.Range().Ptr(),
			})
		}
		return diags
	})
}
