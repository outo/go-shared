package period_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestPeriod(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Period Suite")
}
