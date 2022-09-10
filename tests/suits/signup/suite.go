package signup

import (
	"bytes"
	"context"
	"fmt"
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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	apiServer "apart-deal-api/pkg/api/server"
	authDomain "apart-deal-api/pkg/domain/auth"
	pkgTools "apart-deal-api/pkg/tools"
)

type specContainer struct {
	fx.In

	Echo *echo.Echo
}

var signUpConstModule = fx.Options(
	fx.Supply(&config.Config{
		IsDebug: true,
	}),
	fx.Supply(&dependencies.ApiConfig{
		Port: 37800 + GinkgoParallelProcess(),
	}),
	fx.Provide(apiServer.NewServer),
	fx.Provide(apiServer.NewAuthRouteGroup),
	fx.Provide(user.NewUserRepository),
	fx.Provide(auth.NewSignUpHandler),
	fx.Provide(authDomain.NewSignUpService),
	fx.Invoke(auth.RegisterSignUpRoute),
)

func RegisterSuite(db *mongo.Database) {
	Describe("Sign up", func() {

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
				signUpConstModule,
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

		It("Body is invalid", func() {
			body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(400))
			fmt.Println(rec.Body.String())
		})

		It("Email is invalid", func() {
			body := bytes.NewBuffer([]byte(`{"name":"foo", "email": "foo", "password": "barbaris"}`))
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(400))
			Expect(rec.Body.String()).To(ContainSubstring(`[{"path":"email","message":"must be a valid email address"}]`))
		})

		It("Successful signup", func() {
			body := bytes.NewBuffer([]byte(`{"name":"foo", "email": "foo@gmail.com", "password": "barbaris"}`))
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(200))
			Expect(rec.Body.String()).To(MatchRegexp(`{"token":".+"}`))
		})

		It("Email is occupied by confirmed user", func() {
			_, err := db.Collection("users").InsertOne(ctx, user.User{
				UID:    pkgTools.NewUUID().String(),
				Name:   "Foo",
				Email:  "foo@gmail.com",
				Status: user.StatusConfirmed,
			})
			Expect(err).To(Succeed())

			body := bytes.NewBuffer([]byte(`{"name":"foo", "email": "foo@gmail.com", "password": "barbaris"}`))
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-up", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(409))
			Expect(rec.Body.String()).To(MatchRegexp(`{"message":"Such user already exists"}`))
		})

	})

}
