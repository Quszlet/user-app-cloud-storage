package models

import (
	"time"
)

// * - null
type Directory struct {
	Id               uint   `json:"id"`
	Name             string `gorm:"not null"`
	Size             int    `gorm:"not null"`
	CountFiles       *int
	CountDirectories *int
	DataCreate       time.Time `gorm:"not null"`
	DataChange       *time.Time
	Path             string `gorm:"not null"`
	UserId           int   `gorm:"not null"`
	// Отношение однин ко многим(у директоии может быть много поддиректорий)
	DirectoryId     *int        
	DirectoryUserId *int         
	// У директории могут быть метаданные файлов
	MetadataFiles *[]MetadataFile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
