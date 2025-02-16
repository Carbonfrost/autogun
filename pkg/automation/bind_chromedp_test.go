package automation_test

import (
	"github.com/Carbonfrost/autogun/pkg/automation"
	"github.com/Carbonfrost/autogun/pkg/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bind", func() {

	Describe("bind Task", func() {

		DescribeTable("examples",
			func(cfg config.Task) {
				var (
					output *automation.Automation
					err    error
				)
				bind := func() {
					output, err = automation.Bind(&config.Automation{
						Tasks: []config.Task{
							cfg,
						},
					}, automation.UsingChromedp)
				}

				Expect(bind).NotTo(Panic())
				Expect(err).NotTo(HaveOccurred())
				Expect(output.Tasks).To(HaveLen(1))
				Expect(output.Tasks[0]).NotTo(BeNil())
			},
			Entry("click", new(config.Click)),
			Entry("blur", new(config.Blur)),
			Entry("eval", new(config.Eval)),
			Entry("navigate", new(config.Navigate)),
			Entry("navigate_back", new(config.NavigateBack)),
			Entry("navigate_forward", new(config.NavigateForward)),
			Entry("reload", new(config.Reload)),
			Entry("screenshot", new(config.Screenshot)),
			Entry("sleep", new(config.Sleep)),
			Entry("stop", new(config.Stop)),
			Entry("wait_visible", new(config.WaitVisible)),
		)
	})
})
