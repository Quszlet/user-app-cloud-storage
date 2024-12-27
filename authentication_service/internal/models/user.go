package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type MetadataFile struct {
	Id          int       `json:"id" gorm:"primaryKey"`
	Name        string    `gorm:"not null"`
	Size        int       `gorm:"not null"`
	DataCreate  time.Time `gorm:"not null"`
	DataChange  *time.Time
	Path        string `gorm:"not null"`
	DirectoryId int    `gorm:"type:int(5);not null"` // Связь с Directory
	UserId      int    `gorm:"type:int(5);not null"` // Владелец файла
}

type User struct {
	Id       int    `json:"id" gorm:"primaryKey"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsBanned bool   `json:"is_banned"`
}

type Directory struct {
	Id               uint   `json:"id"`
	Name             string `gorm:"not null"`
	Size             int    `gorm:"not null"`
	CountFiles       *int
	CountDirectories *int
	DataCreate       time.Time `gorm:"not null"`
	DataChange       *time.Time
	Path             string `gorm:"not null"`
	UserId           int    `gorm:"not null"`
	// Отношение однин ко многим(у директоии может быть много поддиректорий)
	DirectoryId     *int
	DirectoryUserId *int
	// У директории могут быть метаданные файлов
	MetadataFiles *[]MetadataFile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Login, validation.Required),
		validation.Field(&u.Email, validation.Required),
		validation.Field(&u.Password, validation.Required))
}
