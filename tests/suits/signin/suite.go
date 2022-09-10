package signin

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"apart-deal-api/dependencies"
	"apart-deal-api/pkg/api/auth"
	"apart-deal-api/pkg/config"
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
)

type specContainer struct {
	fx.In

	Echo *echo.Echo
}

var constModule = fx.Options(
	fx.Supply(&config.Config{
		IsDebug: true,
	}),
	fx.Supply(&dependencies.ApiConfig{
		Port: 37800 + GinkgoParallelProcess(),
	}),
	fx.Provide(apiServer.NewServer),
	fx.Provide(apiServer.NewAuthRouteGroup),
	fx.Provide(user.NewUserRepository),
	fx.Provide(auth.NewAuthenticationService),
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
				fx.Provide(func() auth.TokenStore {
					return auth.NewMockTokenStore(t)
				}),
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
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/sign-in", body)
			req.Header.Add("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			spec.Echo.ServeHTTP(rec, req)

			Expect(rec.Code).To(Equal(400))
			fmt.Println(rec.Body.String())
		})

	})

}
