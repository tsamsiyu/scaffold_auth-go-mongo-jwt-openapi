package tests

import (
	"fmt"
	"testing"

	"apart-deal-api/dependencies"
	"apart-deal-api/tests/suits/signin"
	"apart-deal-api/tests/suits/signup"
	"apart-deal-api/tests/suits/signup_confirm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEverything(t *testing.T) {
	RegisterFailHandler(Fail)
	RegisterTestingT(t)

	BeforeEach(func() {
		fmt.Println(CurrentSpecReport().LeafNodeText)
	})

	dbCfg, err := dependencies.NewDbConfig()
	Expect(err).To(Succeed())

	dbClient, err := dependencies.NewMongoClient(dbCfg)
	Expect(err).To(Succeed())

	db := dependencies.NewMongoDb(dbClient, dbCfg)

	signup.RegisterSuite(db)
	signup_confirm.RegisterSuite(db)
	signin.RegisterSuite(t, db)

	RunSpecs(t, "Everything")
}
