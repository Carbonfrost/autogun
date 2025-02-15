package config // intentional

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

var _ = Describe("validIdentifier", func() {
	DescribeTable("examples",
		func(name string, matcher types.GomegaMatcher) {
			Expect(validIdentifier(name)).To(matcher)
		},
		Entry("nominal", "variable", BeTrue()),
		Entry("allow empty names", "", BeTrue()),
		Entry("valid with underscore", "_underscore", BeTrue()),
		Entry("invalid", "1", BeFalse()),
	)
})
