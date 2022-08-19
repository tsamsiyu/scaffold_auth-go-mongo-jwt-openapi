package auth

import (
	"context"
	"strconv"
	"time"

	"apart-deal-api/pkg/security"
	"apart-deal-api/pkg/store/user"
	"apart-deal-api/pkg/tools"
	"apart-deal-api/pkg/utils"
)

type EmailOccupiedError struct {
	error
}

type SignUpInput struct {
	Name     string
	Email    string
	Password string
}

type SignUpOutput struct {
	Token string
}

type SignUpService struct {
	userRepo user.UserRepository
}

func NewSignUpService(userRepo user.UserRepository) *SignUpService {
	return &SignUpService{
		userRepo: userRepo,
	}
}

func (s *SignUpService) SignUp(ctx context.Context, input SignUpInput) (SignUpOutput, error) {
	passwordHash, err := security.HashPassword(input.Password)
	if err != nil {
		return SignUpOutput{}, err
	}

	token := utils.RandomString(12)
	code := utils.RandomIntBetween(10000, 99999)

	model := user.User{
		UID:          tools.NewUUID().String(),
		Email:        input.Email,
		Name:         input.Name,
		PasswordHash: passwordHash,
		Status:       user.StatusPending,
		CreatedAt:    time.Now(),
		SignUpReq: &user.SignUpRequest{
			Token: token,
			Code:  strconv.Itoa(code),
		},
	}

	if err := s.userRepo.Create(ctx, &model); err != nil {
		return SignUpOutput{}, err
	}

	return SignUpOutput{
		Token: token,
	}, nil
}
