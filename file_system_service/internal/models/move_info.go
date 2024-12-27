package models

import validation "github.com/go-ozzo/ozzo-validation"

type MoveInfo struct {
	MovedObjId int `json:"moved_obj_id"`
	DirectoryId int `json:"directory_id"`
}

func (mi MoveInfo) Validate() error {
	return validation.ValidateStruct(&mi,
		validation.Field(&mi.MovedObjId, validation.Required))
}