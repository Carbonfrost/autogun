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
				Type: "navigate_back",
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

	mappingTaskBlocks = blockMapping[Task]{
		"navigate":         taskMapping(decodeNavigateBlock),
		"navigate_forward": taskMapping(decodeNavigateForwardBlock),
		"navigate_back":    taskMapping(decodeNavigateBackBlock),
		"eval":             taskMapping(decodeEvalBlock),
		"wait_visible":     taskMapping(decodeWaitVisibleBlock),
		"click":            taskMapping(decodeClickBlock),
		"screenshot":       taskMapping(decodeScreenshotBlock),
	}
)

func decodeAutomationBlock(block *hcl.Block) (*Automation, hcl.Diagnostics) {
	f := new(Automation)
	return reduceTask(
		f,
		block,
		supportsDeclRange(&f.DeclRange),
		supportsOptionalLabel(&f.Name, &f.NameRange),
		supportsPartialContentSchema(
			automationBlockSchema,
			appendsTo(&f.Tasks, mappingTaskBlocks),
		),
	)
}
