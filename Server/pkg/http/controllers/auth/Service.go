package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthResult struct {
	Id          string
	Email       string
	Name        string
	AccessToken string
}

type AuthService interface {
	Login(ctx context.Context, email string, password string) (*AuthResult, error)
	SignUp(ctx context.Context, email string, passowrd string, name string) (*AuthResult, error)
}

type authutils struct {
	User          UserRepo
	TokenProvider TokenProvider
}

func NewAuthService(repo UserRepo, token TokenProvider) AuthService {
	return &authutils{repo, token}
}

func (r *authutils) Login(ctx context.Context, email string, password string) (*AuthResult, error) {
	userDetails, userPassword, err := r.User.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	verifyPassword := bcrypt.CompareHashAndPassword([]byte(*userPassword), []byte(password))
	if verifyPassword != nil {
		return nil, errors.New("WRONG PASSWORD")
	}

	accessToken, err := r.TokenProvider.GenerateToken(userDetails.Id.String(), userDetails.Name, userDetails.Email)
	if err != nil {
		return nil, err
	}
	return &AuthResult{
		Id:          userDetails.Id.String(),
		Email:       userDetails.Email,
		Name:        userDetails.Name,
		AccessToken: accessToken,
	}, nil
}

func (r *authutils) SignUp(ctx context.Context, email string, name string, password string) (*AuthResult, error) {
	cost := bcrypt.DefaultCost
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return nil, err
	}
	savedUser, err := r.User.CreateUser(ctx, email, name, string(hashedPassword))
	if err != nil {
		return nil, err
	}
	accessToken, err := r.TokenProvider.GenerateToken(savedUser.Id.String(), email, name)

	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Id:          savedUser.Id.String(),
		Email:       savedUser.Email,
		Name:        savedUser.Name,
		AccessToken: accessToken,
	}, nil
}
