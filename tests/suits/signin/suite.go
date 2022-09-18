package signin

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"apart-deal-api/dependencies"
	"apart-deal-api/pkg/config"
	"apart-deal-api/pkg/security"
	"apart-deal-api/pkg/store/user"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	authHandlers "apart-deal-api/pkg/api/handlers/auth"
	apiServer "apart-deal-api/pkg/api/server"
	pkgTools "apart-deal-api/pkg/tools"
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
		Port:        37800 + GinkgoParallelProcess(),
		TokenSecret: "foobar",
	}),
	fx.Provide(apiServer.NewServer),
	fx.Provide(apiServer.NewAuthRouteGroup),
	fx.Provide(user.NewUserRepository),
	fx.Provide(dependencies.NewAuthenticationService),
	fx.Provide(authHandlers.NewSignInHandler),
	fx.Invoke(authHandlers.RegisterSignInRoute),
)

func RegisterSuite(t *testing.T, db *mongo.Database) {
	Describe("Sign In", func() {

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

			appStartCtx, appStartCancel := context.WithTimeout(context.Background(), time.Second*3)
			defer appStartCancel()

			err = app.Start(appStartCtx)
			Expect(err).To(Succeed())
		})

		AfterEach(func() {
			appStopCtx, appStopCancel := context.WithTimeout(context.Background(), time.Second*3)
			defer appStopCancel()

			err := app.Stop(appStopCtx)
			Expect(err).To(Succeed())

			cancel()
		})

		It("Invalid body", func() {
			body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(400))
		})

		It("Invalid email", func() {
			body := bytes.NewBuffer([]byte(`{"email":"foo@bar","password":"secret"}`))
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(400))
			Expect(rec.Body.String()).To(ContainSubstring(`{"entries":[{"path":"email","message":"must be a valid email address"}],"tag":"validation"}`))
		})

		It("User does not exist", func() {
			body := bytes.NewBuffer([]byte(`{"email":"foo@bar.baz","password":"secret"}`))
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(404))
			Expect(rec.Body.String()).To(ContainSubstring(`{"message":"User not found"}`))
		})

		It("Password is wrong", func() {
			rawPass := "my_secret"
			passHash, err := security.HashPassword(rawPass)
			Expect(err).To(Succeed())

			_, err = db.Collection("users").InsertOne(ctx, user.User{
				UID:          pkgTools.NewUUID().String(),
				Name:         "Foo",
				Email:        "foo@bar.baz",
				PasswordHash: passHash,
				Status:       user.StatusConfirmed,
			})
			Expect(err).To(Succeed())

			body := bytes.NewBuffer([]byte(`{"email":"foo@bar.baz","password":"secret"}`))
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(400))
			Expect(rec.Body.String()).To(ContainSubstring(`{"message":"Password is invalid","tag":"invalid_pass"}`))
		})

		It("Credentials are right", func() {
			rawPass := "my_secret"
			passHash, err := security.HashPassword(rawPass)
			Expect(err).To(Succeed())

			_, err = db.Collection("users").InsertOne(ctx, user.User{
				UID:          pkgTools.NewUUID().String(),
				Name:         "Foo",
				Email:        "foo@bar.baz",
				PasswordHash: passHash,
				Status:       user.StatusConfirmed,
			})
			Expect(err).To(Succeed())

			body := bytes.NewBuffer([]byte(`{"email":"foo@bar.baz","password":"my_secret"}`))
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(200))
			Expect(rec.Body.String()).To(MatchRegexp(`{"token":".*"}`))
		})
	})

}
