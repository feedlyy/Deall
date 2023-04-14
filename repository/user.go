package repository

import (
	"Deall/domain"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) domain.UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Store(ctx context.Context, user domain.Users) error {
	var (
		err error
	)

	_, err = u.db.Collection("users").InsertOne(ctx, user)
	if err != nil {
		logrus.Errorf("Users - Repository|err when store data, err:%v", err)
		return err
	}

	return nil
}
