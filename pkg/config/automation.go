package config

import (
	"github.com/hashicorp/hcl/v2"
)

type Automation struct {
	DeclRange hcl.Range
	NameRange hcl.Range
	Name      string
	Tasks     []Task
}

var (
	automationBlockSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{},
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type: "navigate",
			},
			{
				Type: "navigate_forward",
			},
			{
				Type:       "eval",
				LabelNames: []string{"name"},
			},
			{
				Type: "click",
			},
			{
				Type: "wait_visible",
			},
			{
				Type:       "screenshot",
				LabelNames: []string{"name"},
			},
		},
	}
)

func decodeAutomationBlock(block *hcl.Block) (*Automation, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	f := &Automation{
		Name:      tryLabel(block, 0),
		DeclRange: block.DefRange,
		NameRange: tryLabelRange(block, 0),
	}
	content, _, moreDiags := block.Body.PartialContent(automationBlockSchema)
	diags = append(diags, moreDiags...)

	for _, block := range content.Blocks {
		switch block.Type {

		case "navigate":
			cfg, cfgDiags := decodeNavigateBlock(block)
			diags = append(diags, cfgDiags...)
			if cfg != nil {
				f.Tasks = append(f.Tasks, cfg)
			}

		case "navigate_forward":
			cfg, cfgDiags := decodeNavigateForwardBlock(block)
			diags = append(diags, cfgDiags...)
			if cfg != nil {
				f.Tasks = append(f.Tasks, cfg)
			}

		case "eval":
			cfg, cfgDiags := decodeEvalBlock(block)
			diags = append(diags, cfgDiags...)
			if cfg != nil {
				f.Tasks = append(f.Tasks, cfg)
			}

		case "wait_visible":
			cfg, cfgDiags := decodeWaitVisibleBlock(block)
			diags = append(diags, cfgDiags...)
			if cfg != nil {
				f.Tasks = append(f.Tasks, cfg)
			}

		case "click":
			cfg, cfgDiags := decodeClickBlock(block)
			diags = append(diags, cfgDiags...)
			if cfg != nil {
				f.Tasks = append(f.Tasks, cfg)
			}

		case "screenshot":
			cfg, cfgDiags := decodeScreenshotBlock(block)
			diags = append(diags, cfgDiags...)
			if cfg != nil {
				f.Tasks = append(f.Tasks, cfg)
			}

		default:
			continue
		}
	}

	return f, diags
}
