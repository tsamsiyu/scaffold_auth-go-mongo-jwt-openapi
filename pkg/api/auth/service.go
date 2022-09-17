package auth

import (
	"context"
	"net/http"
	"time"

	"apart-deal-api/pkg/domain/auth"
	"apart-deal-api/pkg/security"
	"apart-deal-api/pkg/utils"

	"github.com/pkg/errors"

	userStore "apart-deal-api/pkg/store/user"

	oas "gitlab.com/apart-deals/openapi/go/api"
)

const (
	TokenLen         int = 12
	TokenExpDuration     = time.Minute * 15
)

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

func (s *AuthenticationService) Take(ctx context.Context, userID string, tokenHash string) (*Token, error) {
	tokens, err := s.tokenStore.FindByKey(ctx, userID)
	if err != nil {
		return nil, err
	}

	var token *Token

	for i, t := range tokens {
		if t.Hash == tokenHash {
			token = &tokens[i]
			break
		}
	}

	return token, nil
}

func (s *AuthenticationService) RefreshToken(ctx context.Context, payload *oas.RefreshAuthToken) (*oas.AuthToken, error) {
	var expiredToken bool
	var newToken *Token

	if err := s.tokenStore.FindForUpdate(ctx, payload.UserId, func(ctx context.Context, tokens []Token) error {
		tokenIndex := -1

		for i, t := range tokens {
			if t.RefreshingHash == payload.RefreshToken && t.Hash == payload.AuthToken {
				tokenIndex = i
				break
			}
		}

		if tokenIndex == -1 {
			return nil
		}

		if !tokens[tokenIndex].Created.Add(TokenExpDuration).Before(time.Now()) {
			expiredToken = true
			return nil
		}

		newToken = &Token{
			UserUID:        payload.UserId,
			Hash:           GenerateTokenHash(),
			RefreshingHash: GenerateTokenRefreshingHash(),
		}

		return s.tokenStore.SetByIndex(ctx, payload.UserId, tokenIndex, *newToken)
	}); err != nil {
		return nil, err
	}

	if expiredToken {
		return nil, &TokenExpiredError{}
	}

	if newToken == nil {
		return nil, &TokenDoesNotExistError{}
	}

	return &oas.AuthToken{
		UserId:       payload.UserId,
		Token:        newToken.Hash,
		RefreshToken: newToken.RefreshingHash,
		ExpiresAt:    newToken.Created.Add(TokenExpDuration),
	}, nil
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
		return nil, &UserNotConfirmedError{error: errors.New("Not authorized")}
	}

	if ok := security.CheckPasswordHash(payload.Password, user.PasswordHash); !ok {
		return nil, &InvalidPasswordError{error: errors.New("Invalid password")}
	}

	token := Token{
		UserUID:        user.UID,
		Hash:           GenerateTokenHash(),
		RefreshingHash: GenerateTokenRefreshingHash(),
		Created:        time.Now(),
	}

	if err := s.tokenStore.Push(ctx, user.UID, token); err != nil {
		return nil, err
	}

	return &oas.AuthToken{
		UserId:       token.UserUID,
		Token:        token.Hash,
		RefreshToken: token.RefreshingHash,
		ExpiresAt:    token.Created.Add(TokenExpDuration),
	}, nil
}

func GenerateTokenHash() string {
	return utils.RandomString(TokenLen)
}
func GenerateTokenRefreshingHash() string {
	return utils.RandomString(TokenLen)
}

func ExtractTokenFromHttpHeaders(header http.Header) string {
	h := header.Get("Authorization")

	return h
}
