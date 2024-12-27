package models

import validation "github.com/go-ozzo/ozzo-validation"

// Как нибудь вынести в общий модуль
type User struct {
	Id          int   `json:"id" gorm:"primaryKey"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	IsBanned    bool   `json:"is_banned"`
	Directories []Directory `gorm:"constraint:OnDelete:CASCADE;"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Login, validation.Required),
		validation.Field(&u.Email, validation.Required),
		validation.Field(&u.Password, validation.Required))
}
