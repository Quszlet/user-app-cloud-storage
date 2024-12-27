package repository

import (
	"github.com/Quszlet/user_service/internal/models"
	"gorm.io/gorm"
)

type User interface {
	Create(user models.User) (int, error)
	Update(user models.User) error
	GetById(userId int) (models.User, error)
	GetByLogin(login string) (models.User, error)
	GetAll() ([]models.User, error)
	Delete(userId int) error
	Migration()
}

type Repository struct {
	User
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{User: NewUserMySql(db)}
}
