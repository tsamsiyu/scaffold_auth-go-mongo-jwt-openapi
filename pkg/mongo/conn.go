package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MainDB *mongo.Database

type TLSConfig struct {
	CaCrt     string
	ClientKey string
	ClientCrt string
}

type DbRef struct {
	URI string
	TLS *TLSConfig
}

func NewClient(ctx context.Context, ref *DbRef) (*mongo.Client, error) {
	opts := options.Client().
		ApplyURI(ref.URI).
		SetTimeout(time.Second * 5)

	if ref.TLS != nil {
		caPool := x509.NewCertPool()

		caCert, err := os.ReadFile(ref.TLS.CaCrt)
		if err != nil {
			return nil, err
		}

		caPem, _ := pem.Decode(caCert)

		caX509, err := x509.ParseCertificate(caPem.Bytes)
		if err != nil {
			return nil, err
		}

		clientTlsCert, err := tls.LoadX509KeyPair(ref.TLS.ClientCrt, ref.TLS.ClientKey)
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

func ProvideDatabase(client *mongo.Client, dbname string) MainDB {
	return client.Database(dbname)
}
