package auth

import (
	"context"
	"net/http"
	"time"

	"apart-deal-api/pkg/domain/auth"
	"apart-deal-api/pkg/utils"

	"github.com/pkg/errors"

	userStore "apart-deal-api/pkg/store/user"

	oas "gitlab.com/apart-deals/openapi/go/api"
)

const (
	TokenLen         int = 12
	TokenExpDuration     = time.Minute * 15
)

type UserNotConfirmed struct {
	error
}

type AuthenticationService struct {
	tokenStore TokenStore
	userRepo   userStore.UserRepository
}

func NewAuthenticationService(
	tokenStore TokenStore,
	userRepo userStore.UserRepository,
) *AuthenticationService {
	return &AuthenticationService{
		tokenStore: tokenStore,
		userRepo:   userRepo,
	}
}

func (s *AuthenticationService) Auth(ctx context.Context, payload *oas.SignIn) (*oas.AuthToken, error) {
	user, err := s.userRepo.FindByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &auth.UserNotFound{}
	}

	if user.Status != userStore.StatusConfirmed {
		return nil, &UserNotConfirmed{error: errors.New("Not authorized")}
	}

	tokenModel := &Token{
		UserUID:        user.UID,
		Hash:           utils.RandomString(TokenLen),
		RefreshingHash: utils.RandomString(TokenLen),
		Created:        time.Now(),
	}

	if err := s.tokenStore.Create(ctx, tokenModel, TokenExpDuration); err != nil {
		return nil, err
	}

	return &oas.AuthToken{
		Token:        tokenModel.Hash,
		RefreshToken: tokenModel.RefreshingHash,
		ExpiresAt:    tokenModel.Created.Add(TokenExpDuration),
	}, nil
}

func ExtractTokenFromHttpHeaders(header http.Header) string {
	h := header.Get("Authorization")

	return h
}
