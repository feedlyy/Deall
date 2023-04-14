package repository

import (
	"Deall/domain"
	"Deall/helpers"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	_, err = u.db.Collection(helpers.UsersCollection).InsertOne(ctx, user)
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

	cursor, err = u.db.Collection(helpers.UsersCollection).Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"password": 0}))
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

func (u *userRepository) GetByUsername(ctx context.Context, usr string) (domain.Users, error) {
	var (
		err error
		res domain.Users
	)

	err = u.db.Collection(helpers.UsersCollection).FindOne(ctx, bson.M{"username": usr}).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.Error("User - Repository|no document found")
			return domain.Users{}, err
		}
		logrus.Errorf("User - Repository|err when get by username, err:%v", err)
		return domain.Users{}, err
	}

	return res, nil
}

func (u *userRepository) Delete(ctx context.Context, id string) error {
	var (
		err    error
		objID  primitive.ObjectID
		delRes *mongo.DeleteResult
	)

	// turn string id into objID
	objID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		logrus.Errorf("User - Repository|err when generate objID, err:%v", err)
		return err
	}

	delRes, err = u.db.Collection(helpers.UsersCollection).DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		logrus.Errorf("User - Repository|err when delete user by username, err:%v", err)
		return err
	}

	if delRes.DeletedCount == 0 {
		err = mongo.ErrNoDocuments
		logrus.Error("User - Repository|no document found")
		return err
	}

	return nil
}
