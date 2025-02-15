package automation_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAutomation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Automation Suite")
}
