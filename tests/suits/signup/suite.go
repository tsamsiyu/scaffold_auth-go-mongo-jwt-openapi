package signup

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"apart-deal-api/pkg/api/handlers/auth"
	"apart-deal-api/pkg/store/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/fx"

	. "apart-deal-api/tests/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	apiServer "apart-deal-api/pkg/api/server"
	authDomain "apart-deal-api/pkg/domain/auth"
	pkgTools "apart-deal-api/pkg/tools"
)

var runSpec = BuildApiSpecRunner(
	fx.Provide(apiServer.NewAuthRouteGroup),
	fx.Provide(user.NewUserRepository),
	fx.Provide(auth.NewSignUpHandler),
	fx.Provide(authDomain.NewSignUpService),
	fx.Invoke(auth.RegisterSignUpRoute),
)

var _ = Describe("Sign up", func() {
	BeforeEach(func() {
		_, err := Shared().Db.Collection("users").DeleteMany(context.Background(), bson.M{})
		Expect(err).To(Succeed())
	})

	It("Body is invalid", func() {
		runSpec(func(apiClient *http.Client) SpecRunner {
			return func(ctx context.Context) {
				body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
				resp, err := apiClient.Post("/api/v1/auth/sign-up", "application/json", body)
				Expect(err).To(Succeed())
				Expect(resp.StatusCode).To(Equal(400))
			}
		})
	})

	It("Email is invalid", func() {
		runSpec(func(apiClient *http.Client) SpecRunner {
			return func(ctx context.Context) {
				body := bytes.NewBuffer([]byte(`{"name":"foo", "email": "foo", "password": "barbaris"}`))
				resp, err := apiClient.Post("/api/v1/auth/sign-up", "application/json", body)
				Expect(err).To(Succeed())

				defer resp.Body.Close()
				respBody, _ := io.ReadAll(resp.Body)

				Expect(err).To(Succeed())
				Expect(resp.StatusCode).To(Equal(400))
				Expect(string(respBody)).To(ContainSubstring(`[{"path":"email","message":"must be a valid email address"}]`))
			}
		})
	})

	It("Successful signup", func() {
		runSpec(func(apiClient *http.Client) SpecRunner {
			return func(ctx context.Context) {
				body := bytes.NewBuffer([]byte(`{"name":"foo", "email": "foo@gmail.com", "password": "barbaris"}`))
				resp, err := apiClient.Post("/api/v1/auth/sign-up", "application/json", body)

				defer resp.Body.Close()
				respBody, _ := io.ReadAll(resp.Body)

				Expect(err).To(Succeed())
				Expect(resp.StatusCode).To(Equal(200))
				Expect(string(respBody)).To(MatchRegexp(`{"token":".+"}`))
			}
		})
	})

	It("Email is occupied by confirmed user", func() {
		runSpec(func(apiClient *http.Client) SpecRunner {
			return func(ctx context.Context) {
				_, err := Shared().Db.Collection("users").InsertOne(ctx, user.User{
					UID:    pkgTools.NewUUID().String(),
					Name:   "Foo",
					Email:  "foo@gmail.com",
					Status: user.StatusConfirmed,
				})
				Expect(err).To(Succeed())

				body := bytes.NewBuffer([]byte(`{"name":"foo", "email": "foo@gmail.com", "password": "barbaris"}`))
				resp, err := apiClient.Post("/api/v1/auth/sign-up", "application/json", body)

				defer resp.Body.Close()
				respBody, _ := io.ReadAll(resp.Body)

				Expect(err).To(Succeed())
				Expect(resp.StatusCode).To(Equal(409))
				Expect(string(respBody)).To(MatchRegexp(`{"message":"Such user already exists"}`))
			}
		})
	})

})
