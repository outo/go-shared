package period_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	//"github.com/pkg/profile"
	"testing"
)

func TestPeriod(t *testing.T) {
	//defer profile.Start(profile.CPUProfile).Stop()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Period Suite")
}
