package service

import "github.com/Quszlet/authentication_service/internal/models"

type Authentication interface {
	SignIn(Username, Password string) (string, error)
	GetUser(Username, Password string) (*models.User, error)
	GenerateJWTToken(UserId int, Username string) (string, error)
}

type Service struct {
	Authentication
}

func NewService(au *AuthService) *Service {
	return &Service{Authentication: au}
}