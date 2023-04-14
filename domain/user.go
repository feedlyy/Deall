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
	GetByUsername(ctx context.Context, usr string) (Users, error)
	Delete(ctx context.Context, id string) error
}

type UserService interface {
	Register(ctx context.Context, user Users) error
	GetUsers(ctx context.Context) ([]Users, error)
	GetUser(ctx context.Context, usr string) (Users, error)
	DeleteUser(ctx context.Context, id string) error
}
