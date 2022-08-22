package common

import (
	"context"
	"fmt"
	"net/http"

	"apart-deal-api/dependencies"
	"apart-deal-api/pkg/config"
	"apart-deal-api/pkg/mongo/schema"
	"apart-deal-api/tests/tools"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	apiServer "apart-deal-api/pkg/api/server"
)

type SpecRunner func()

var DbClient *mongo.Client
var Db *mongo.Database
var Logger *zap.Logger

var _ = BeforeSuite(func() {
	Logger = dependencies.LoggerFromEnv()

	dbCfg, err := dependencies.NewDbConfig()
	Expect(err).To(Succeed())

	DbClient, err = dependencies.NewMongoClient(dbCfg)
	Expect(err).To(Succeed())

	Db = dependencies.NewMongoDb(DbClient, dbCfg)

	err = schema.UsersMigrations(context.Background(), Db)
	Expect(err).To(Succeed())
})

var _ = AfterSuite(func() {
	_ = DbClient.Disconnect(context.Background())
})

func BuildApiSpecRunner(additionalProviders ...fx.Option) func(specRunnerProvider interface{}) {
	return func(specRunnerProvider interface{}) {
		apiProviders := []fx.Option{
			fx.Supply(Logger),
			fx.Supply(Db),
			fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: logger}
			}),
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
			fx.Provide(specRunnerProvider),
			fx.Invoke(func(lc fx.Lifecycle, apiRunFn dependencies.ApiRunFn, e *echo.Echo, specFn SpecRunner, shutdowner fx.Shutdowner) {
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
		}

		apiProviders = append(apiProviders, additionalProviders...)

		app := fx.New(apiProviders...)

		app.Run()

		if err := app.Err(); err != nil {
			Logger.With(zap.Error(err)).Warn("Spec runner returned error")
		}
	}
}
