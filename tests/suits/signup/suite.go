package signup

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"apart-deal-api/dependencies"
	"apart-deal-api/pkg/api/handlers/auth"
	"apart-deal-api/pkg/config"
	"apart-deal-api/pkg/store/user"
	"apart-deal-api/tests/tools"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	apiServer "apart-deal-api/pkg/api/server"
	authDomain "apart-deal-api/pkg/domain/auth"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type specFnType func()

var dbClient *mongo.Client
var db *mongo.Database
var logger *zap.Logger

var _ = BeforeSuite(func() {
	logger = dependencies.LoggerFromEnv()

	dbCfg, err := dependencies.NewDbConfig()
	Expect(err).To(Succeed())

	dbClient, err = dependencies.NewMongoClient(dbCfg)
	Expect(err).To(Succeed())

	db = dependencies.NewMongoDb(dbClient, dbCfg)
})

var _ = AfterSuite(func() {
	_ = dbClient.Disconnect(context.Background())
})

func runSpec(specFnProvider interface{}) {
	app := fx.New(
		fx.Supply(logger),
		fx.Supply(db),
		fx.Supply(&config.Config{
			IsDebug: true,
		}),
		fx.Supply(&dependencies.ApiConfig{
			Port: 37800 + GinkgoParallelProcess(),
		}),
		fx.Provide(func() *http.Client {
			return &http.Client{
				Transport: &tools.MyHttpTransport{
					Host:          fmt.Sprintf("localhost:%d", 37800+GinkgoParallelProcess()),
					BaseTransport: http.DefaultTransport,
				},
			}
		}),
		fx.Provide(dependencies.NewApiRunFn),
		fx.Provide(apiServer.NewServer),
		fx.Provide(apiServer.NewAuthRouteGroup),
		fx.Provide(user.NewUserRepository),
		fx.Provide(auth.NewSignUpHandler),
		fx.Provide(authDomain.NewSignUpService),
		fx.Provide(specFnProvider),
		fx.Invoke(auth.RegisterSignUpRoute),
		fx.Invoke(func(lc fx.Lifecycle, apiRunFn dependencies.ApiRunFn, e *echo.Echo, specFn specFnType, shutdowner fx.Shutdowner) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := apiRunFn(context.Background())
					if err != nil {
						return err
					}

					defer GinkgoRecover()
					defer shutdowner.Shutdown()

					specFn()

					return nil
				},
				OnStop: func(ctx context.Context) error {
					_ = e.Shutdown(ctx)

					return nil
				},
			})
		}),
	)

	app.Run()
}

var _ = Describe("My Tests", func() {

	It("Test 0", func() {
		runSpec(func(apiClient *http.Client) specFnType {
			return func() {
				body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
				resp, err := apiClient.Post("api/v1/auth/sign-up", "application/json", body)
				Expect(err).To(Succeed())
				Expect(resp.StatusCode).To(Equal(400))
			}
		})
	})

})
