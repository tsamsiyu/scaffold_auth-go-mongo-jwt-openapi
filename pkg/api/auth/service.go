package auth

import (
	"context"
	"time"

	"apart-deal-api/pkg/security"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"

	userStore "apart-deal-api/pkg/store/user"

	oas "gitlab.com/apart-deals/openapi/go/api"
)

type TokenPayload struct {
	UserID string
	Email  string
}

const (
	TokenExpDuration = time.Minute * 15
)

type AuthenticationService struct {
	tokenSecret string
	userRepo    userStore.UserRepository
}

func NewAuthenticationService(
	tokenSecret string,
	userRepo userStore.UserRepository,
) *AuthenticationService {
	return &AuthenticationService{
		tokenSecret: tokenSecret,
		userRepo:    userRepo,
	}
}

func (s *AuthenticationService) Sign(payload TokenPayload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    payload.UserID,
		"userEmail": payload.Email,
		"nbf":       time.Now(),
		"exp":       time.Now().Add(TokenExpDuration),
	})

	tokenString, err := token.SignedString([]byte(s.tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthenticationService) Verify(tokenString string) (*TokenPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &TokenInvalidError{}
		}

		return s.tokenSecret, nil
	})
	if err != nil {
		return nil, &TokenInvalidError{}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, &TokenInvalidError{}
	}

	if err := claims.Valid(); err != nil {
		return nil, &TokenInvalidError{}
	}

	return &TokenPayload{
		UserID: claims["userID"].(string),
		Email:  claims["userEmail"].(string),
	}, nil
}

func (s *AuthenticationService) FindUser(ctx context.Context, payload *oas.SignIn) (*userStore.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &NoSuchUserError{}
	}

	if user.Status != userStore.StatusConfirmed {
		return nil, &UserNotConfirmedError{error: errors.New("Not authorized")}
	}

	if ok := security.CheckPasswordHash(payload.Password, user.PasswordHash); !ok {
		return nil, &InvalidPasswordError{error: errors.New("Invalid password")}
	}

	return user, nil
}

func (s *AuthenticationService) Auth(ctx context.Context, payload *oas.SignIn) (string, error) {
	user, err := s.FindUser(ctx, payload)
	if err != nil {
		return "", err
	}

	tokenString, err := s.Sign(TokenPayload{
		UserID: user.UID,
	})
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
