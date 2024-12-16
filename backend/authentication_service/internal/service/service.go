package service

type Authentication interface {
	SignIn(urerName, Password string) (string, error)
}

type Service struct {
	Authentication
}

func NewService() *Service {
	return &Service{}
}