package auth

//go:generate mockery --dir=. --name=TokenStore --inpackage --filename=token_store_mock.go --with-expecter

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Token struct {
	UserUID        string    `json:"userUID"`
	Hash           string    `json:"hash"`
	RefreshingHash string    `json:"refreshingHash"`
	Created        time.Time `json:"created"`
}

type TokenStore interface {
	FindByHash(ctx context.Context, hash string) (*Token, error)
	Create(ctx context.Context, token *Token, expiration time.Duration) error
	DeleteByHash(ctx context.Context, hash string) error
}

func NewTokenStore(client *redis.Client) TokenStore {
	return &redisTokenStorage{client: client}
}

type redisTokenStorage struct {
	client *redis.Client
}

func (s *redisTokenStorage) FindByHash(ctx context.Context, hash string) (*Token, error) {
	cmd := s.client.Get(ctx, hash)
	if err := cmd.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	var token Token

	cmdBytes, err := cmd.Bytes()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if err := json.Unmarshal(cmdBytes, &token); err != nil {
		return nil, errors.WithStack(err)
	}

	return &token, nil
}

func (s *redisTokenStorage) Create(ctx context.Context, token *Token, expiration time.Duration) error {
	serialized, err := json.Marshal(token)
	if err != nil {
		return errors.WithStack(err)
	}

	if token.Hash == "" || token.UserUID == "" {
		return errors.New("Either hash and userID must be present")
	}

	cmd := s.client.Set(ctx, token.Hash, serialized, expiration)
	if err := cmd.Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *redisTokenStorage) DeleteByHash(ctx context.Context, hash string) error {
	cmd := s.client.Del(ctx, hash)
	if err := cmd.Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
