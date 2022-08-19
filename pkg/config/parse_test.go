package config_test

import (
	"os"

	"github.com/Carbonfrost/autogun/pkg/config"
	"github.com/spf13/afero"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/types"
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
						"URL": Equal("https://example.com"),
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
			})),

			Entry("wait_visible", "wait_visible.autog", MatchElementsWithIndex(IndexIdentity, IgnoreExtras, Elements{
				"1": And(
					BeAssignableToTypeOf(&config.WaitVisible{}),
					PointTo(MatchFields(IgnoreExtras, Fields{
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
