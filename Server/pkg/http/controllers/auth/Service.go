package auth

import "context"

type AuthResult struct {
	Id    string
	Email string
	Name  string
}

type AuthService interface {
	Login(ctx context.Context, email string, password string) (*AuthResult, error)
	SignUp(ctx context.Context, email string, passowrd string, name string, emailAppPassword *string) *AuthService.error
}
