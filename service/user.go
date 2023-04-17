package service

import (
	"Deall/domain"
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type userService struct {
	userRepo  domain.UserRepository
	tokenRepo domain.TokenRepository
}

func NewUserService(u domain.UserRepository, t domain.TokenRepository) domain.UserService {
	return &userService{
		userRepo:  u,
		tokenRepo: t,
	}
}

func (u *userService) Register(ctx context.Context, user domain.Users) error {
	var (
		err error
		pwd string
	)

	// bcrypt the password
	pwd, err = domain.HashPassword(user.Password)
	if err != nil {
		logrus.Errorf("User - Service|err when hash password, err:%v", err)
		return err
	}
	user.Password = pwd
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

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

func (u *userService) UpdateUser(ctx context.Context, user domain.Users) error {
	var (
		err error
		pwd string
	)

	if user.Password != "" {
		// bcrypt the password
		pwd, err = domain.HashPassword(user.Password)
		if err != nil {
			logrus.Errorf("User - Service|err when hash password, err:%v", err)
			return err
		}
	}
	user.Password = pwd
	user.UpdatedAt = time.Now()

	return u.userRepo.Update(ctx, user)
}

func (u *userService) Authentication(ctx context.Context, username, password string) (string, error) {
	var (
		token = domain.GenerateRandomUUID()
		err   error
		usr   domain.Users
	)

	// check if valid username and password
	usr, err = u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	err = domain.ValidatePassword(usr.Password, password)
	if err != nil {
		logrus.Errorf("User - Service|err when validate password, err:%v", err)
		return "", err
	}

	// create token
	err = u.tokenRepo.Store(ctx, domain.Tokens{
		UserID:     usr.ID,
		Token:      token,
		Expiration: time.Now().Add(time.Hour),
		CreatedAt:  time.Now(),
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
