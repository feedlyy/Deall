package service

import (
	"Deall/domain"
	"context"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(u domain.UserRepository) domain.UserService {
	return &userService{
		userRepo: u,
	}
}

func (u *userService) Register(ctx context.Context, user domain.Users) error {
	return u.userRepo.Store(ctx, user)
}
