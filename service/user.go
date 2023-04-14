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

func (u *userService) GetUsers(ctx context.Context) ([]domain.Users, error) {
	return u.userRepo.Fetch(ctx)
}

func (u *userService) GetUser(ctx context.Context, usr string) (domain.Users, error) {
	return u.userRepo.GetByUsername(ctx, usr)
}

func (u *userService) DeleteUser(ctx context.Context, id string) error {
	return u.userRepo.Delete(ctx, id)
}
