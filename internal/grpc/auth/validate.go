package auth

import (
	"errors"

	"github.com/go-playground/validator/v10"
	ssov1 "github.com/kuromii5/proto-auth/gen/go/sso"
)

var (
	ErrInvalidEmail  = errors.New("invalid email address")
	ErrShortPassword = errors.New("min password length is 8")
	ErrLongEmail     = errors.New("max email length is 254")
	ErrLongPassword  = errors.New("max password length is 64")
	ErrRequired      = errors.New("this field is required")
)

type RegisterRequest struct {
	Email    string `validate:"required,email,max=254"`
	Password string `validate:"required,min=8,max=64"`
}

func (s *serverAPI) validateRegisterRequest(req *ssov1.RegisterRequest) error {
	v := RegisterRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	if err := s.validate.Struct(v); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, ve := range validationErrors {
				switch ve.Field() {
				case "Email":
					if ve.Tag() == "max" {
						return ErrLongEmail
					}
					if ve.Tag() == "email" {
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

func (s *serverAPI) validateLoginRequest(req *ssov1.LoginRequest) error {
	v := LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	if err := s.validate.Struct(v); err != nil {
		return errors.New("invalid credentials")
	}

	return nil
}
