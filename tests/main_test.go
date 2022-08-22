package tests

import (
	"testing"

	_ "apart-deal-api/tests/common"
	_ "apart-deal-api/tests/suits/signup"
	_ "apart-deal-api/tests/suits/signup_confirm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEverything(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Everything")
}
