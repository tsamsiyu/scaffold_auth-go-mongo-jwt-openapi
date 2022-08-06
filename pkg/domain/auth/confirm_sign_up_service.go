package auth

import (
	"context"

	"apart-deal-api/pkg/store/user"
)

type ConfirmationCodeMismatchError struct {
	error
}

type ConfirmSignUpInput struct {
	Token string
	Code  string
}

type ConfirmSignUpService struct {
	userRepo user.UserRepository
}

func NewConfirmSignUpService(userRepo user.UserRepository) *ConfirmSignUpService {
	return &ConfirmSignUpService{
		userRepo: userRepo,
	}
}

func (s *ConfirmSignUpService) Confirm(ctx context.Context, input ConfirmSignUpInput) error {
	userModel, err := s.userRepo.FindBySignUpReqToken(ctx, input.Token)
	if err != nil {
		return err
	}

	if userModel.SignUpReq.Code != input.Code {
		return &ConfirmationCodeMismatchError{}
	}

	if err := s.userRepo.ConfirmAndDeleteSignUpReq(ctx, userModel.UID); err != nil {
		return err
	}

	return nil
}
