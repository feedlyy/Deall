package domain

import (
	"context"
	"time"
)

type Users struct {
	Username  string    `bson:"username"`
	Password  string    `bson:"password" json:"-"`
	Role      string    `bson:"role"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type UserRepository interface {
	Store(ctx context.Context, user Users) error
	Fetch(ctx context.Context) ([]Users, error)
}

type UserService interface {
	Register(ctx context.Context, user Users) error
	GetUSers(ctx context.Context) ([]Users, error)
}
