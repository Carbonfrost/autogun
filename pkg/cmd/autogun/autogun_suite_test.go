package autogun_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAutogun(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Autogun Suite")
}
