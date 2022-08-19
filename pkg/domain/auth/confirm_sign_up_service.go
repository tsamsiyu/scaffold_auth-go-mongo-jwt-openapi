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

	if userModel == nil {
		return &UserNotFound{}
	}

	if userModel.SignUpReq.Code != input.Code {
		return &ConfirmationCodeMismatchError{}
	}

	confirmed, err := s.userRepo.ConfirmAndDeleteSignUpReq(ctx, userModel.UID)
	if err != nil {
		return err
	}

	if !confirmed {
		return &CouldNotConfirmError{}
	}

	return nil
}
