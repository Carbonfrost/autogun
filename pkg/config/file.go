package config

import (
	"github.com/hashicorp/hcl/v2"
)

type File struct {
	Automations []*Automation
}

var (
	fileSchema = &hcl.BodySchema{
		Attributes: []hcl.AttributeSchema{},
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "automation",
				LabelNames: []string{"name"},
			},
		},
	}
)

func decodeFile(body hcl.Body) (*File, hcl.Diagnostics) {
	var diags hcl.Diagnostics
	f := &File{}

	content, contentDiags := body.Content(fileSchema)
	diags = append(diags, contentDiags...)

	for _, block := range content.Blocks {
		switch block.Type {

		case "automation":
			cfg, cfgDiags := decodeAutomationBlock(block)
			diags = append(diags, cfgDiags...)
			if cfg != nil {
				f.Automations = append(f.Automations, cfg)
			}

		default:
			continue
		}
	}

	return f, diags
}
