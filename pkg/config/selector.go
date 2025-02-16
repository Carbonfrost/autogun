package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

type Selector struct {
	DeclRange hcl.Range
	Target    string
	By        SelectorBy
	On        SelectorOn
}

type SelectorBy string
type SelectorOn string

const (
	BySearch   SelectorBy = "SEARCH"
	ByJSPath   SelectorBy = "JS_PATH"
	ByID       SelectorBy = "ID"
	ByQuery    SelectorBy = "QUERY"
	ByQueryAll SelectorBy = "QUERY_ALL"
)

const (
	OnReady      SelectorOn = "READY"
	OnVisible    SelectorOn = "VISIBLE"
	OnNotVisible SelectorOn = "NOT_VISIBLE"
	OnEnabled    SelectorOn = "ENABLED"
	OnSelected   SelectorOn = "SELECTED"
	OnNotPresent SelectorOn = "NOT_PRESENT"
)

var (
	selectorBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{
			{Name: "target"},
			{Name: "by"},
			{Name: "on"},
		},
		Blocks: []hcl.BlockHeaderSchema{},
	}
)

func decodeSelectorBlock(block *hcl.Block) (*Selector, hcl.Diagnostics) {
	f := new(Selector)
	return reduce(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsPartialContentSchema(
			selectorBlockSchema,
			withAttribute("target", &f.Target),
			withAttr("by", func(attr *hcl.Attribute) hcl.Diagnostics {
				var in string
				diags := gohcl.DecodeExpression(attr.Expr, nil, &in)
				moreDiags := parseSelectorBy(in, &f.By, &attr.Range)
				return append(diags, moreDiags...)
			}),
			withAttr("on", func(attr *hcl.Attribute) hcl.Diagnostics {
				var in string
				diags := gohcl.DecodeExpression(attr.Expr, nil, &in)
				moreDiags := parseSelectorOn(in, &f.On, &attr.Range)
				return append(diags, moreDiags...)
			}),
		),
	)
}

func parseSelectorBy(s string, dst *SelectorBy, subject *hcl.Range) hcl.Diagnostics {
	switch s {
	case "SEARCH", "search":
		*dst = BySearch
		return nil
	case "JS_PATH", "js_path":
		*dst = ByJSPath
		return nil
	case "ID", "id":
		*dst = ByID
		return nil
	case "QUERY", "query":
		*dst = ByQuery
		return nil
	case "QUERY_ALL", "query_all":
		*dst = ByQueryAll
		return nil
	}
	return hcl.Diagnostics{
		diagInvalidValue(s, "by", subject),
	}
}

func parseSelectorOn(s string, dst *SelectorOn, subject *hcl.Range) hcl.Diagnostics {
	switch s {
	case "READY", "ready":
		*dst = OnReady
		return nil
	case "VISIBLE", "visible":
		*dst = OnVisible
		return nil
	case "NOT_VISIBLE", "not_visible":
		*dst = OnNotVisible
		return nil
	case "ENABLED", "enabled":
		*dst = OnEnabled
		return nil
	case "SELECTED", "selected":
		*dst = OnSelected
		return nil
	case "NOT_PRESENT", "not_present":
		*dst = OnNotPresent
		return nil
	}
	return hcl.Diagnostics{
		diagInvalidValue(s, "on", subject),
	}
}

func supportsSelectorBlocks(sels *[]*Selector, opts **Options) partialContentMapper {
	return func(content *hcl.BodyContent) hcl.Diagnostics {
		var diags hcl.Diagnostics

		for _, block := range content.Blocks {
			switch block.Type {
			case "selector":
				cfg, cfgDiags := decodeSelectorBlock(block)
				if cfg != nil {
					*sels = append(*sels, cfg)
				}
				diags = append(diags, cfgDiags...)

			case "options":
				cfg, cfgDiags := decodeOptionsBlock(block)
				if cfg != nil {
					*opts = cfg
				}
				diags = append(diags, cfgDiags...)

			default:
				continue
			}
		}
		return diags
	}
}
