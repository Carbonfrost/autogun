package config_test

import (
	"os"
	"time"

	"github.com/Carbonfrost/autogun/pkg/config"
	"github.com/hashicorp/hcl/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
	"github.com/spf13/afero"
)

var _ = Describe("LoadConfigFile", func() {

	Describe("parse Task", func() {
		DescribeTable("examples",
			func(hclFile string, expected types.GomegaMatcher) {
				res, err := validExample(hclFile)
				Expect(err).NotTo(HaveOccurred())
				Expect(res.Automations[0].Tasks).To(expected)
			},

			Entry("navigate", "eval.autog", MatchElementsWithIndex(IndexIdentity, IgnoreExtras, Elements{
				"0": And(
					BeAssignableToTypeOf(&config.Navigate{}),
					PointTo(MatchFields(IgnoreExtras, Fields{
						"URL": WithTransform(toString, Equal("https://example.com")),
					}))),
			})),

			Entry("evaluate", "eval.autog", MatchElementsWithIndex(IndexIdentity, IgnoreExtras, Elements{
				"1": And(
					BeAssignableToTypeOf(&config.Eval{}),
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Name":   Equal("output"),
						"Script": Equal("1"),
					}))),
			})),

			Entry("click", "click.autog", MatchElementsWithIndex(IndexIdentity, IgnoreExtras, Elements{
				"1": And(
					BeAssignableToTypeOf(&config.Click{}),
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Selector": Equal("#olive"),
					}))),
				"2": And(
					BeAssignableToTypeOf(&config.Click{}),
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Selectors": MatchElementsWithIndex(IndexIdentity, IgnoreExtras, Elements{
							"0": And(
								BeAssignableToTypeOf(&config.Selector{}),
								PointTo(MatchFields(IgnoreExtras, Fields{
									"Target": Equal("#raspberry"),
								}))),
						}),
					}))),
				"3": And(
					BeAssignableToTypeOf(&config.Click{}),
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Selectors": MatchElementsWithIndex(IndexIdentity, IgnoreExtras, Elements{
							"0": And(
								BeAssignableToTypeOf(&config.Selector{}),
								PointTo(MatchFields(IgnoreExtras, Fields{
									"By": Equal(config.BySearch),
									"On": Equal(config.OnVisible),
								}))),
						}),
					}))),
			})),

			Entry("wait_visible", "wait_visible.autog", MatchElementsWithIndex(IndexIdentity, IgnoreExtras, Elements{
				"1": And(
					BeAssignableToTypeOf(&config.WaitVisible{}),
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Selector": Equal("#aubergine"),
						"Options": PointTo(MatchFields(IgnoreExtras, Fields{
							"AtLeast":       PointTo(Equal(1)),
							"RetryInterval": PointTo(Equal(5 * time.Second)),
						})),
					}))),
			})),

			Entry("screenshot", "screenshot.autog", MatchElementsWithIndex(IndexIdentity, IgnoreExtras, Elements{
				"1": And(
					BeAssignableToTypeOf(&config.Screenshot{}),
					PointTo(MatchFields(IgnoreExtras, Fields{
						"Name":     Equal("label.png"),
						"Selector": Equal("#aubergine"),
					}))),
			})),
		)
	})
})

func validExample(hclFile string) (*config.File, error) {
	appFS := afero.NewMemMapFs()
	appFS.MkdirAll(".weyoun/", 0755)
	data, err := os.ReadFile("testdata/valid-examples/" + hclFile)
	Expect(err).NotTo(HaveOccurred())

	afero.WriteFile(appFS, ".weyoun/site.hcl", data, 0644)

	p := config.NewParser(afero.NewIOFS(appFS))
	return p.LoadConfigFile(".weyoun/site.hcl")
}

func toString(v interface{}) interface{} {
	d, _ := v.(hcl.Expression).Value(nil)
	return d.AsString()
}
