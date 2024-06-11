package auth

import (
	ssov1 "github.com/kuromii5/proto-auth/gen/go/sso"
)

// Validation structs
type RegisterRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}
type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}
type IsAdminRequest struct {
	UserID int64 `validate:"required,gt=0"`
}

// Convert and validate functions
func (s *serverAPI) validateRegisterRequest(req *ssov1.RegisterRequest) (*RegisterRequest, error) {
	v := &RegisterRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	if err := s.validate.Struct(v); err != nil {
		return nil, err
	}
	return v, nil
}
func (s *serverAPI) validateLoginRequest(req *ssov1.LoginRequest) (*LoginRequest, error) {
	v := &LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	if err := s.validate.Struct(v); err != nil {
		return nil, err
	}
	return v, nil
}
func (s *serverAPI) validateIsAdminRequest(req *ssov1.IsAdminRequest) (*IsAdminRequest, error) {
	v := &IsAdminRequest{
		UserID: req.GetUserId(),
	}
	if err := s.validate.Struct(v); err != nil {
		return nil, err
	}
	return v, nil
}
