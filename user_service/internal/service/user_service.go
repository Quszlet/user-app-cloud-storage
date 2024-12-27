package service

import (
	"github.com/Quszlet/user_service/internal/models"
	"github.com/Quszlet/user_service/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(r *repository.Repository) *UserService {
	return &UserService{repo: r.User}
}

func (us *UserService) Create(user models.User) (int, error) {
	return us.repo.Create(user)
}

func (us *UserService) Update(user models.User) error {
	return us.repo.Update(user)
}

func (us *UserService) GetById(userId int) (models.User, error) {
	return us.repo.GetById(userId)
}

func (us *UserService) GetByLogin(login string) (models.User, error) {
	return us.repo.GetByLogin(login)
}

func (us *UserService) GetAll() ([]models.User, error) {
	return us.repo.GetAll()
}

func (us *UserService) Delete(userId int) error {
	return us.repo.Delete(userId)
}