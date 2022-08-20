package signup

import (
	"apart-deal-api/pkg/store/user"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"

	"apart-deal-api/dependencies"
	"apart-deal-api/pkg/api/handlers/auth"
	"apart-deal-api/pkg/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/zap"

	apiServer "apart-deal-api/pkg/api/server"
	authDomain "apart-deal-api/pkg/domain/auth"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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

func runSpec(specFn interface{}) {
	app := fx.New(
		fx.Supply(logger),
		fx.Supply(db),
		fx.Supply(&config.Config{
			IsDebug: true,
		}),
		fx.Supply(&dependencies.ApiConfig{
			Port: 37800 + GinkgoParallelProcess(),
		}),
		fx.Provide(dependencies.NewApiRunFn),
		fx.Provide(apiServer.NewServer),
		fx.Provide(apiServer.NewAuthRouteGroup),
		fx.Provide(user.NewUserRepository),
		fx.Provide(auth.NewSignUpHandler),
		fx.Provide(authDomain.NewSignUpService),
		fx.Invoke(auth.RegisterSignUpRoute),
		fx.Invoke(func(lc fx.Lifecycle, apiRunFn dependencies.ApiRunFn, e *echo.Echo) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return apiRunFn(context.Background())
				},
				OnStop: func(ctx context.Context) error {
					_ = e.Shutdown(ctx)

					return nil
				},
			})
		}),
		fx.Invoke(specFn),
	)

	app.Run()
}

var _ = Describe("My Tests", func() {

	It("Test 0", func() {
		runSpec(func(shutdowner fx.Shutdowner) error {
			fmt.Println("OK 0")
			return shutdowner.Shutdown()
		})
	})

	It("Test 1", func() {
		runSpec(func(shutdowner fx.Shutdowner) error {
			fmt.Println("OK 1")
			return shutdowner.Shutdown()
		})
	})

})
