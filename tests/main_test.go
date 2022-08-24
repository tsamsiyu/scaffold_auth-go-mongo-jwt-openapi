package tests

import (
	"context"
	"testing"

	"apart-deal-api/tests/common"

	_ "apart-deal-api/tests/common"
	_ "apart-deal-api/tests/suits/signup"
	_ "apart-deal-api/tests/suits/signup_confirm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEverything(t *testing.T) {
	RegisterFailHandler(Fail)
	RegisterTestingT(t)

	SynchronizedBeforeSuite(func() []byte {
		return []byte("")
	}, func(bytes []byte) {
		common.InitSharedDeps(context.Background())
	})

	SynchronizedAfterSuite(func() {
		common.CleanupSharedDeps()
	}, func() {
	})

	RunSpecs(t, "Everything")
}
