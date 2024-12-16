package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type User struct {
	Id        uint   `json:"id"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Is_banned bool   `json:"is_banned"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Login, validation.Required),
		validation.Field(&u.Email, validation.Required),
		validation.Field(&u.Password, validation.Required))
}
