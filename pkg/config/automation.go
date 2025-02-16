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
				Type: "blur",
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
			{
				Type: "sleep",
			},
			{
				Type: "reload",
			},
			{
				Type: "stop",
			},
		},
	}

	mappingTaskBlocks = blockMapping[Task]{
		"blur":             taskMapping(decodeBlurBlock),
		"click":            taskMapping(decodeClickBlock),
		"eval":             taskMapping(decodeEvalBlock),
		"navigate":         taskMapping(decodeNavigateBlock),
		"navigate_back":    taskMapping(decodeNavigateBackBlock),
		"navigate_forward": taskMapping(decodeNavigateForwardBlock),
		"screenshot":       taskMapping(decodeScreenshotBlock),
		"reload":           taskMapping(decodeReloadBlock),
		"sleep":            taskMapping(decodeSleepBlock),
		"stop":             taskMapping(decodeStopBlock),
		"wait_visible":     taskMapping(decodeWaitVisibleBlock),
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
