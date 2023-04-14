package repository

import (
	"Deall/domain"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (u *userRepository) Fetch(ctx context.Context) ([]domain.Users, error) {
	var (
		err    error
		res    []domain.Users
		cursor *mongo.Cursor
	)

	cursor, err = u.db.Collection("users").Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"password": 0}))
	if err != nil {
		logrus.Errorf("Users - Repository|err when fetch all data, err:%v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var usr domain.Users

		if err = cursor.Decode(&usr); err != nil {
			logrus.Errorf("User - Repository|err when decode cursor, err:%v", err)
			return nil, err
		}
		res = append(res, usr)
	}

	return res, nil
}
