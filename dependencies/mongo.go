package dependencies

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"os"
	"time"

	"apart-deal-api/pkg/mongo/schema"

	"github.com/Netflix/go-env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/fx"
)

type DbConfig struct {
	DbUri       string `env:"MONGO_URI,required=true"`
	DbName      string `env:"MONGO_DOMAIN_DB,required=true"`
	DbCaCrt     string `env:"MONGO_CA_CRT"`
	DbClientCrt string `env:"MONGO_CLIENT_CRT"`
	DbClientKey string `env:"MONGO_CLIENT_KEY"`
}

type TLSConfig struct {
	CaCrt     string
	ClientKey string
	ClientCrt string
}

type DbRef struct {
	URI string
	TLS *TLSConfig
}

func NewMongoClient(cfg *DbConfig) (*mongo.Client, error) {
	opts := options.Client().
		ApplyURI(cfg.DbUri).
		SetTimeout(time.Second * 5)

	if cfg.DbClientKey != "" {
		caPool := x509.NewCertPool()

		caCert, err := os.ReadFile(cfg.DbCaCrt)
		if err != nil {
			return nil, err
		}

		caPem, _ := pem.Decode(caCert)

		caX509, err := x509.ParseCertificate(caPem.Bytes)
		if err != nil {
			return nil, err
		}

		clientTlsCert, err := tls.LoadX509KeyPair(cfg.DbClientCrt, cfg.DbClientKey)
		if err != nil {
			return nil, err
		}

		caPool.AddCert(caX509)

		opts = opts.SetTLSConfig(&tls.Config{
			RootCAs:            caPool,
			Certificates:       []tls.Certificate{clientTlsCert},
			InsecureSkipVerify: true,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client, err := mongo.Connect(ctx, opts)

	readPref, err := readpref.New(readpref.PrimaryPreferredMode)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readPref); err != nil {
		return nil, err
	}

	return client, err
}

func NewDbConfig() (*DbConfig, error) {
	var cfg DbConfig

	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func NewMongoDb(client *mongo.Client, cfg *DbConfig) *mongo.Database {
	return client.Database(cfg.DbName)
}

var DbModule = fx.Module("Mongo",
	fx.Provide(
		NewDbConfig,
		NewMongoClient,
		NewMongoDb,
	),
	fx.Invoke(func(lc fx.Lifecycle, client *mongo.Client, db *mongo.Database) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := schema.Migrate(ctx, db); err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				_ = client.Disconnect(ctx)

				return nil
			},
		})
	}),
)
