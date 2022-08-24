package signup_confirm

import (
	"bytes"
	"net/http"

	"apart-deal-api/pkg/api/handlers/auth"
	"apart-deal-api/pkg/store/user"

	"go.uber.org/fx"

	apiServer "apart-deal-api/pkg/api/server"
	authDomain "apart-deal-api/pkg/domain/auth"

	. "apart-deal-api/tests/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var runSpec = BuildApiSpecRunner(
	fx.Provide(apiServer.NewAuthRouteGroup),
	fx.Provide(user.NewUserRepository),
	fx.Provide(auth.NewSignUpConfirmHandler),
	fx.Provide(authDomain.NewConfirmSignUpService),
	fx.Invoke(auth.RegisterSignUpConfirmRoute),
)

var _ = Describe("Sign Up Confirmation", func() {

	It("Invalid body", func() {
		runSpec(func(apiClient *http.Client) SpecRunner {
			return func() {
				body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
				resp, err := apiClient.Post("/api/v1/auth/sign-up-confirm", "application/json", body)
				Expect(err).To(Succeed())
				Expect(resp.StatusCode).To(Equal(400))
			}
		})
	})

})
