package auth

import (
	"errors"

	"github.com/go-playground/validator/v10"
	sso "github.com/kuromii5/miku-notes-auth/generated"
)

var (
	ErrInvalidEmail  = errors.New("invalid email address")
	ErrShortPassword = errors.New("min password length is 8")
	ErrLongPassword  = errors.New("max password length is 64")
	ErrRequired      = errors.New("this field is required")
)

type RegisterRequest struct {
	Email    string `validate:"required,email,max=254"`
	Password string `validate:"required,min=8,max=64"`
}

func validateRegisterRequest(req *sso.RegisterRequest) error {
	validate := validator.New()

	v := RegisterRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	if err := validate.Struct(v); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, ve := range validationErrors {
				switch ve.Field() {
				case "Email":
					if ve.Tag() == "email" || ve.Tag() == "max" {
						return ErrInvalidEmail
					}
					return ErrRequired

				case "Password":
					if ve.Tag() == "min" {
						return ErrShortPassword
					} else if ve.Tag() == "max" {
						return ErrLongPassword
					}
					return ErrRequired

				default:
					return errors.New(ve.Error())
				}
			}
		}

		return err
	}

	return nil
}

type LoginRequest struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}

func validateLoginRequest(req *sso.LoginRequest) error {
	validate := validator.New()

	v := LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	if err := validate.Struct(v); err != nil {
		return errors.New("invalid credentials")
	}

	return nil
}
