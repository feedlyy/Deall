package domain

import (
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (u Users) ValidateUser() error {
	var err error

	if u.Username == "" {
		return errors.New("username cannot be empty")
	} else if u.Password == "" {
		return errors.New("password cannot be empty")
	} else if u.Role == "" {
		return errors.New("role cannot be empty")
	}

	if err = u.ValidateRole(); err != nil {
		return err
	}

	return nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidatePassword(hashedPassword string, plainPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		return err
	}

	return nil
}

func (u Users) ValidateRole() error {
	if u.Role != "Admin" && u.Role != "User" {
		return errors.New("role are not exists, only accepted Admin or User")
	}
	return nil
}

func GenerateRandomUUID() string {
	// Generate a new UUID
	id := uuid.New()

	return id.String()
}

func LocalLocation() (*time.Location, error) {
	var (
		loc *time.Location
		err error
	)

	loc, err = time.LoadLocation("Asia/Jakarta")
	if err != nil {
		logrus.Errorf("Helper|Err when get  %v", err)
		return nil, err
	}

	return loc, nil
}

func IsValidContains(arr []string, s string) bool {
	for _, val := range arr {
		if s == val {
			return true
		}
	}
	return false
}
