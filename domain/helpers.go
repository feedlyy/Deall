package domain

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
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

func (u Users) ValidateRole() error {
	if u.Role != "Admin" && u.Role != "User" {
		return errors.New("role are not exists, only accepted Admin or User")
	}
	return nil
}
