package users

import (
	"errors"
	"net/mail"
	"strings"
	"unicode"
)

type UserSignup struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserSignup) checkFields() error {
	if len(strings.Trim(u.Username, " ")) == 0 {
		return errors.New("Username is a required field.")
	}

	if len(u.Username) < 6 {
		return errors.New("Username must be at least 6 characters long.")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("Please enter a valid email address.")
	}

	if err := isValid(u.Password); err != nil {
		return err
	}

	return nil
}

func isValid(password string) error {
	var (
		hasMinLen = false
		hasUpper  = false
		hasLower  = false
		hasNumber = false
	)
	if len(password) >= 7 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper {
		return errors.New("Password must contain at least one uppercase letter.")
	}

	if !hasMinLen {
		return errors.New("Password must be at least 7 characters long")
	}

	if !hasLower {
		return errors.New("Password must contain at least one lowercase letter.")
	}

	if !hasNumber {
		return errors.New("Password must contain at least one number.")
	}

	return nil
}

type UserSignin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserSignin) checkFields() error {
	if len(strings.Trim(u.Email, " ")) == 0 {
		return errors.New("Please enter a valid email address.")
	}

	if len(strings.Trim(u.Password, " ")) == 0 {
		return errors.New("Please enter a valid password.")
	}

	return nil
}
