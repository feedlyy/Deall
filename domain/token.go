package domain

import (
	"context"
	"time"
)

type Tokens struct {
	UserID     string    `json:"user_id,omitempty"`
	Token      string    `json:"token,omitempty"`
	Expiration time.Time `json:"-"`
	CreatedAt  time.Time `json:"-"`
}

type TokenRepository interface {
	Store(ctx context.Context, token Tokens) error
	GetByToken(ctx context.Context, token string) (Tokens, error)
}
