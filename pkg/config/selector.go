// Copyright 2025 The Autogun Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package config

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/hcl/v2"
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

func (s *Selector) setOn(n SelectorOn) {
	s.On = n
}

func (s *Selector) setBy(n SelectorBy) {
	s.By = n
}

func decodeSelectorBlock(block *hcl.Block) (*Selector, hcl.Diagnostics) {
	f := new(Selector)
	return reduce(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsPartialContentSchema(
			selectorBlockSchema,
			withAttribute("target", &f.Target),
			withAttributeParser("by", f.setBy, parseSelectorBy),
			withAttributeParser("on", f.setOn, parseSelectorOn),
		),
	)
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func parseSelectorBy(s string) (result SelectorBy, err error) {
	switch s {
	case "SEARCH", "search":
		return BySearch, nil
	case "JS_PATH", "js_path":
		return ByJSPath, nil
	case "ID", "id":
		return ByID, nil
	case "QUERY", "query":
		return ByQuery, nil
	case "QUERY_ALL", "query_all":
		return ByQueryAll, nil
	}
	err = fmt.Errorf("value %q is not a valid value", s)
	return
}

func parseSelectorOn(s string) (result SelectorOn, err error) {
	switch s {
	case "READY", "ready":
		return OnReady, nil
	case "VISIBLE", "visible":
		return OnVisible, nil
	case "NOT_VISIBLE", "not_visible":
		return OnNotVisible, nil
	case "ENABLED", "enabled":
		return OnEnabled, nil
	case "SELECTED", "selected":
		return OnSelected, nil
	case "NOT_PRESENT", "not_present":
		return OnNotPresent, nil
	}
	err = fmt.Errorf("value %q is not a valid value", s)
	return
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
