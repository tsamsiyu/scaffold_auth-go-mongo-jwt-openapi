package auth

//go:generate mockery --dir=. --name=TokenStore --inpackage --filename=token_store_mock.go --with-expecter

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

const (
	ctxCmdableKey = "redis_cmdable"
)

func WithCmdable(ctx context.Context, cmdable redis.Cmdable) context.Context {
	return context.WithValue(ctx, ctxCmdableKey, cmdable)
}

func CmdableFromCtx(ctx context.Context) redis.Cmdable {
	v := ctx.Value(ctxCmdableKey)
	if v == nil {
		return nil
	}

	cmdable, ok := v.(redis.Cmdable)
	if !ok {
		return nil
	}

	return cmdable
}

type Token struct {
	UserUID        string    `json:"userUID"`
	Hash           string    `json:"hash"`
	RefreshingHash string    `json:"refreshingHash"`
	Created        time.Time `json:"created"`
}

type TokenStore interface {
	FindByKey(ctx context.Context, key string) ([]Token, error)
	FindForUpdate(ctx context.Context, key string, fn func(context.Context, []Token) error) error
	Push(ctx context.Context, key string, token Token) error
	SetByIndex(ctx context.Context, key string, index int, t Token) error
}

func NewTokenStore(client *redis.Client) TokenStore {
	return &redisTokenStorage{client: client}
}

type redisTokenStorage struct {
	client *redis.Client
}

func (s *redisTokenStorage) FindByKey(ctx context.Context, key string) ([]Token, error) {
	var cmdable redis.Cmdable

	pipe := CmdableFromCtx(ctx)
	if pipe != nil {
		cmdable = pipe
	} else {
		cmdable = s.client
	}

	cmd := cmdable.LRange(ctx, key, 0, -1)
	if err := cmd.Err(); err != nil {
		return nil, errors.WithStack(err)
	}

	val, err := cmd.Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tokens := make([]Token, 0)

	for _, item := range val {
		var token Token

		if err := json.Unmarshal([]byte(item), &token); err != nil {
			return nil, errors.WithStack(err)
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (s *redisTokenStorage) FindForUpdate(ctx context.Context, key string, fn func(context.Context, []Token) error) error {
	if err := s.client.Watch(ctx, func(tx *redis.Tx) error {
		tokens, err := s.FindByKey(WithCmdable(ctx, tx), key)
		if err != nil {
			return err
		}

		if _, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			if err := fn(WithCmdable(ctx, pipe), tokens); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
		return nil
	}, key); err != nil {
		return err
	}

	return nil
}

func (s *redisTokenStorage) Push(ctx context.Context, key string, token Token) error {
	if token.Hash == "" || token.UserUID == "" {
		return errors.New("Either hash and userID must be present")
	}

	serialized, err := json.Marshal(token)
	if err != nil {
		return errors.WithStack(err)
	}

	var cmdable redis.Cmdable

	pipe := CmdableFromCtx(ctx)
	if pipe != nil {
		cmdable = pipe
	} else {
		cmdable = s.client
	}

	cmd := cmdable.LPush(ctx, key, serialized)
	if err := cmd.Err(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *redisTokenStorage) SetByIndex(ctx context.Context, key string, index int, t Token) error {
	var cmdable redis.Cmdable

	pipe := CmdableFromCtx(ctx)
	if pipe != nil {
		cmdable = pipe
	} else {
		cmdable = s.client
	}

	tJson, err := json.Marshal(t)
	if err != nil {
		return err
	}

	cmd := cmdable.LSet(ctx, key, int64(index), tJson)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}
