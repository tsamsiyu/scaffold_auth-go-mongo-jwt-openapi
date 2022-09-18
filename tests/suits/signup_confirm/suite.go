package signup_confirm

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"

	"apart-deal-api/dependencies"
	"apart-deal-api/pkg/api/handlers/auth"
	"apart-deal-api/pkg/config"
	"apart-deal-api/pkg/store/user"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	apiServer "apart-deal-api/pkg/api/server"
	authDomain "apart-deal-api/pkg/domain/auth"
	pkgTools "apart-deal-api/pkg/tools"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type specContainer struct {
	fx.In

	Echo *echo.Echo
}

var constModule = fx.Options(
	fx.NopLogger,
	fx.Supply(&config.Config{
		IsDebug: true,
	}),
	fx.Supply(&dependencies.ApiConfig{
		Port: 37800 + GinkgoParallelProcess(),
	}),
	fx.Provide(apiServer.NewServer),
	fx.Provide(apiServer.NewAuthRouteGroup),
	fx.Provide(user.NewUserRepository),
	fx.Provide(auth.NewSignUpConfirmHandler),
	fx.Provide(authDomain.NewConfirmSignUpService),
	fx.Invoke(auth.RegisterSignUpConfirmRoute),
)

func RegisterSuite(db *mongo.Database) {
	Describe("Sign Up Confirmation", func() {
		loggerLvl := zap.NewAtomicLevelAt(zap.ErrorLevel)
		logger := dependencies.NewLogger(&loggerLvl)

		var (
			ctx    context.Context
			cancel context.CancelFunc
			app    *fx.App
			spec   *specContainer
		)

		BeforeEach(func() {
			ctx, cancel = context.WithCancel(context.Background())

			_, err := db.Collection("users").DeleteMany(ctx, bson.M{})
			Expect(err).To(Succeed())

			app = fx.New(
				fx.Supply(logger),
				fx.Supply(db),
				constModule,
				fx.Invoke(func(s specContainer) {
					spec = &s
				}),
			)

			err = app.Start(context.Background())
			Expect(err).To(Succeed())
		})

		AfterEach(func() {
			err := app.Stop(context.Background())
			Expect(err).To(Succeed())

			cancel()
		})

		It("Invalid body", func() {
			body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up-confirm", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(400))
		})

		It("User with given token does not exist", func() {
			body := bytes.NewBuffer([]byte(`{"code":"123","token":"qwe"}`))
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up-confirm", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(404))
		})

		It("User's code is wrong", func() {
			_, err := db.Collection("users").InsertOne(ctx, user.User{
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
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up-confirm", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(400))
			Expect(rec.Body.String()).To(ContainSubstring("code_mismatch"))
		})

		It("User's code is correct", func() {
			userUID := pkgTools.NewUUID().String()
			_, err := db.Collection("users").InsertOne(ctx, user.User{
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
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up-confirm", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			res := db.Collection("users").FindOne(ctx, bson.D{
				{"_id", userUID},
			})
			Expect(res.Err()).To(Succeed())

			var found user.User

			err = res.Decode(&found)
			Expect(err).To(Succeed())

			Expect(found.Status).To(Equal(user.StatusConfirmed))
			Expect(found.SignUpReq).To(BeNil())
		})

		It("User is already confirmed", func() {
			_, err := db.Collection("users").InsertOne(ctx, user.User{
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
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up-confirm", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(400))
			Expect(rec.Body.String()).To(ContainSubstring("unconfirmable"))
		})

	})

}
