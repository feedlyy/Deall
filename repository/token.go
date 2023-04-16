package repository

import (
	"Deall/domain"
	"Deall/helpers"
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type tokenRepository struct {
	db *mongo.Database
}

func NewTokenRepository(db *mongo.Database) domain.TokenRepository {
	return &tokenRepository{db: db}
}

func (t *tokenRepository) Store(ctx context.Context, token domain.Tokens) error {
	var (
		err error
	)

	_, err = t.db.Collection(helpers.TokenCollection).InsertOne(ctx, token)
	if err != nil {
		logrus.Errorf("Token - Repository|err when store data, err:%v", err)
		return err
	}

	return nil
}

func (t *tokenRepository) GetByToken(ctx context.Context, token string) (domain.Tokens, error) {
	var (
		err    error
		res    domain.Tokens
		filter = bson.M{}
	)

	filter = bson.M{
		"$and": bson.A{
			bson.M{"token": token},
			bson.M{"expiration": bson.M{"$gte": time.Now()}},
		},
	}

	err = t.db.Collection(helpers.TokenCollection).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.Error("Token - Repository|no document found")
			return domain.Tokens{}, err
		}
		logrus.Errorf("Token - Repository|err when get by username, err:%v", err)
		return domain.Tokens{}, err
	}

	return res, nil
}
