package models

import "time"

type MetadataFile struct {
	Id          int      `json:"id" gorm:"primaryKey"`
	Name        string    `gorm:"not null"`
	Size        int       `gorm:"not null"`
	DataCreate  time.Time `gorm:"not null"`
	DataChange  *time.Time
	Path        string    `gorm:"not null"`
	DirectoryId int       `gorm:"type:int(5);not null"`            // Связь с Directory
    UserId      int       `gorm:"type:int(5);not null"`            // Владелец файла
}
