package signup_confirm

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"apart-deal-api/pkg/api/handlers/auth"
	"apart-deal-api/pkg/store/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/fx"

	apiServer "apart-deal-api/pkg/api/server"
	authDomain "apart-deal-api/pkg/domain/auth"
	pkgTools "apart-deal-api/pkg/tools"

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
	BeforeEach(func() {
		_, err := Shared().Db.Collection("users").DeleteMany(context.Background(), bson.M{})
		Expect(err).To(Succeed())
	})

	It("Invalid body", func() {
		runSpec(func(apiClient *http.Client) SpecRunner {
			return func(ctx context.Context) {
				body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
				resp, err := apiClient.Post("/api/v1/auth/sign-up-confirm", "application/json", body)
				Expect(err).To(Succeed())
				Expect(resp.StatusCode).To(Equal(400))
			}
		})
	})

	It("User with given token does not exist", func() {
		runSpec(func(apiClient *http.Client) SpecRunner {
			return func(ctx context.Context) {
				body := bytes.NewBuffer([]byte(`{"code":"123","token":"qwe"}`))
				resp, err := apiClient.Post("/api/v1/auth/sign-up-confirm", "application/json", body)
				Expect(err).To(Succeed())
				Expect(resp.StatusCode).To(Equal(404))
			}
		})
	})

	It("User's code is wrong", func() {
		runSpec(func(apiClient *http.Client) SpecRunner {
			return func(ctx context.Context) {
				_, err := Shared().Db.Collection("users").InsertOne(ctx, user.User{
					UID:    pkgTools.NewUUID().String(),
					Name:   "Foo",
					Email:  "foo@gmail.com",
					Status: user.StatusConfirmed,
					SignUpReq: &user.SignUpRequest{
						Code:       "228",
						Token:      "qwe",
						NotifiedAt: nil,
					},
				})
				Expect(err).To(Succeed())

				body := bytes.NewBuffer([]byte(`{"code":"123","token":"qwe"}`))
				resp, err := apiClient.Post("/api/v1/auth/sign-up-confirm", "application/json", body)
				Expect(err).To(Succeed())

				respBody, err := io.ReadAll(resp.Body)
				Expect(err).To(Succeed())

				Expect(resp.StatusCode).To(Equal(400))
				Expect(respBody).To(ContainSubstring("code_mismatch"))
			}
		})
	})

	It("User's code is correct", func() {
		runSpec(func(apiClient *http.Client) SpecRunner {
			return func(ctx context.Context) {
				userUID := pkgTools.NewUUID().String()
				_, err := Shared().Db.Collection("users").InsertOne(ctx, user.User{
					UID:    userUID,
					Name:   "Foo",
					Email:  "foo@gmail.com",
					Status: user.StatusPending,
					SignUpReq: &user.SignUpRequest{
						Code:       "123",
						Token:      "qwe",
						NotifiedAt: nil,
					},
				})
				Expect(err).To(Succeed())

				body := bytes.NewBuffer([]byte(`{"code":"123","token":"qwe"}`))
				resp, err := apiClient.Post("/api/v1/auth/sign-up-confirm", "application/json", body)
				Expect(err).To(Succeed())
				Expect(resp.StatusCode).To(Equal(204))

				res := Shared().Db.Collection("users").FindOne(ctx, bson.D{
					{"_id", userUID},
				})
				Expect(res.Err()).To(Succeed())

				var found user.User

				err = res.Decode(&found)
				Expect(err).To(Succeed())

				Expect(found.Status).To(Equal(user.StatusConfirmed))
				Expect(found.SignUpReq).To(BeNil())
			}
		})
	})

	It("User is already confirmed", func() {
		runSpec(func(apiClient *http.Client) SpecRunner {
			return func(ctx context.Context) {
				_, err := Shared().Db.Collection("users").InsertOne(ctx, user.User{
					UID:    pkgTools.NewUUID().String(),
					Name:   "Foo",
					Email:  "foo@gmail.com",
					Status: user.StatusConfirmed,
					SignUpReq: &user.SignUpRequest{
						Code:       "123",
						Token:      "qwe",
						NotifiedAt: nil,
					},
				})
				Expect(err).To(Succeed())

				body := bytes.NewBuffer([]byte(`{"code":"123","token":"qwe"}`))
				resp, err := apiClient.Post("/api/v1/auth/sign-up-confirm", "application/json", body)
				Expect(err).To(Succeed())

				respBody, err := io.ReadAll(resp.Body)
				Expect(err).To(Succeed())

				Expect(resp.StatusCode).To(Equal(400))
				Expect(respBody).To(ContainSubstring("unconfirmable"))
			}
		})
	})

})
