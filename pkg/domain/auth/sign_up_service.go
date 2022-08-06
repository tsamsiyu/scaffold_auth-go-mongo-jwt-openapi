package auth

import (
	"context"
	"strconv"
	"time"

	"apart-deal-api/pkg/mongo/db"
	"apart-deal-api/pkg/security"
	"apart-deal-api/pkg/store/user"
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

type SignUpService struct {
	userRepo user.UserRepository
}

func NewSignUpService(userRepo user.UserRepository) *SignUpService {
	return &SignUpService{
		userRepo: userRepo,
	}
}

func (s *SignUpService) SignUp(ctx context.Context, input SignUpInput) error {
	existingUser, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return &EmailOccupiedError{}
	}

	passwordHash, err := security.HashPassword(input.Password)
	if err != nil {
		return err
	}

	token := utils.RandomString(12)
	code := utils.RandomIntBetween(10000, 99999)

	model := user.User{
		UID:          db.NewUUID().String(),
		Email:        input.Email,
		Name:         input.Name,
		PasswordHash: passwordHash,
		Status:       user.Pending,
		SignUpReq: &user.SignUpRequest{
			Token:     token,
			Code:      strconv.Itoa(code),
			CreatedAt: time.Now(),
		},
	}

	if err := s.userRepo.Create(ctx, &model); err != nil {
		return err
	}

	return nil
}
