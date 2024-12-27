package repository

import (
	"fmt"

	"github.com/Quszlet/user_service/internal/models"
	"gorm.io/gorm"
)

type UserMySql struct {
	db *gorm.DB
}

func NewUserMySql(db *gorm.DB) *UserMySql {
	return &UserMySql{db: db}
}

func (up *UserMySql) Migration() {
	up.db.AutoMigrate(&models.User{})

	var count int64
    up.db.Model(&models.User{}).Count(&count)

    if count == 0 {
		users := []models.User{
			{Login: "admin", Password: "admin123", Email: "admin@example.com", IsBanned: false},
			{Login: "user1", Password: "password1", Email: "user1@example.com", IsBanned: false},
			{Login: "user2", Password: "password2", Email: "user2@example.com", IsBanned: true},
			{Login: "moderator", Password: "mod123", Email: "moderator@example.com", IsBanned: false},
			{Login: "guest", Password: "guest123", Email: "guest@example.com", IsBanned: false},
		}
		up.db.Create(&users)
		fmt.Println("Users added in context migration")
	}
}

func (up *UserMySql) Create(user models.User) (int, error) {
	result := up.db.Create(&user);
	if result.Error != nil {
		return 0, result.Error
	}

	return int(user.Id), nil
}

func (up *UserMySql) Update(user models.User) error {
	up.db.Save(user)
	return nil
}

func (up *UserMySql) GetById(userId int) (models.User, error) {
	var user models.User
	result := up.db.First(&user, userId)
	if result.Error != nil {
		return models.User{}, nil
	}
	return user, nil
}

func (up *UserMySql) GetByLogin(login string) (models.User, error) {
	var user models.User
	result := up.db.Where("login = ?", login).First(&user)
	if result.Error != nil {
		return models.User{}, nil
	}
	return user, nil
}

func (up *UserMySql) GetAll() ([]models.User, error) {
	var users []models.User
	result := up.db.Find(&users)
	if result.Error != nil {
		return []models.User{}, nil
	}
	return users, nil
}

func (up *UserMySql) Delete(userId int) error {
	result := up.db.Delete(&models.User{}, userId)
	if result.Error != nil {
		return result.Error
	}
	
	return nil
}