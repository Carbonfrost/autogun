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
	var diags hcl.Diagnostics
	f := &Selector{
		DeclRange: block.DefRange,
	}

	content, _, moreDiags := block.Body.PartialContent(selectorBlockSchema)
	diags = append(diags, moreDiags...)

	if attr, ok := content.Attributes["target"]; ok {
		moreDiags := gohcl.DecodeExpression(attr.Expr, nil, &f.Target)
		diags = append(diags, moreDiags...)
	}

	if attr, ok := content.Attributes["by"]; ok {
		var in string
		moreDiags := gohcl.DecodeExpression(attr.Expr, nil, &in)
		diags = append(diags, moreDiags...)

		moreDiags = parseSelectorBy(in, &f.By, &attr.Range)
		diags = append(diags, moreDiags...)
	}

	if attr, ok := content.Attributes["on"]; ok {
		var in string
		moreDiags := gohcl.DecodeExpression(attr.Expr, nil, &in)
		diags = append(diags, moreDiags...)

		moreDiags = parseSelectorOn(in, &f.On, &attr.Range)
		diags = append(diags, moreDiags...)
	}

	return f, diags
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

func andSelectorQueryAttributes(content *hcl.BodyContent, dst selectorAction) hcl.Diagnostics {
	var diags hcl.Diagnostics

	if attr, ok := content.Attributes["selector"]; ok {
		var dstSelector string
		moreDiags := gohcl.DecodeExpression(attr.Expr, nil, &dstSelector)
		diags = append(diags, moreDiags...)
		dst.setSelector(dstSelector)
	}

	for _, block := range content.Blocks {
		switch block.Type {

		case "selector":
			cfg, cfgDiags := decodeSelectorBlock(block)
			diags = append(diags, cfgDiags...)
			if cfg != nil {
				dst.addSelector(cfg)
			}
		case "options":
			cfg, cfgDiags := decodeOptionsBlock(block)
			diags = append(diags, cfgDiags...)
			if cfg != nil {
				dst.setOptions(cfg)
			}

		default:
			continue
		}
	}

	return diags
}
