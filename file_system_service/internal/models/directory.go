package models

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"gorm.io/gorm"
)

type Directory struct {
	Id               int        `json:"id"`
	Name             string     `json:"name" gorm:"not null;default:home"`
	Size             int        `json:"size" gorm:"not null"`
	CountFiles       *int       `json:"count_files"`
	CountDirectories *int       `json:"count_directories"`
	DataCreate       time.Time  `json:"data_create" gorm:"not null"`
	DataChange       *time.Time `json:"data_change"`
	Path             string     `json:"path" gorm:"not null;default:/home"`
	UserId           int        `json:"user_id" gorm:"not null"`
	// Отношение однин ко многим(у директоии может быть много поддиректорий)
	DirectoryId     *int         `json:"directory_id"`
	DirectoryUserId *int         `json:"directory_user_id"`
	Directories     *[]Directory `gorm:"constraint:OnDelete:CASCADE;"`
	// // У директории могут быть метаданные файлов
	MetadataFiles *[]MetadataFile `gorm:"constraint:OnDelete:CASCADE;"`
}

func (d Directory) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required),
		validation.Field(&d.UserId, validation.Required))
}

func (d *Directory) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Printf("Начало работы триггера(Директория %s) AfterCreate\n", d.Name)
	fmt.Println(d)

	// Домашняя директория
	if d.DirectoryId == nil {
		count := 0
		d.CountDirectories = &count
		return
	}

	parentDirectoryId := d.DirectoryId
	directory := &Directory{}

	tx.Where("id = ?", parentDirectoryId).First(&directory)
	if d.DirectoryId != nil {
		d.Path = directory.Path + "/" + d.Name
	}

	directory = d

	iter := 1
	oldCountDir := 0
	for directory.Id != 0 {
		fmt.Printf("Итерация - %d\n", iter)
		iter++

		newDirectoryCount := 0
		if directory.CountDirectories != nil {
			newDirectoryCount = *directory.CountDirectories + 1 + oldCountDir
		}
		directory.CountDirectories = &newDirectoryCount
		tx.Save(&directory)

		fmt.Printf("Dir name - %s. countDir - %d\n", directory.Name, *directory.CountDirectories)

		if directory.DirectoryId != nil {
			newDirectory := &Directory{}
			parentDirectoryId = directory.DirectoryId

			tx.Where("id = ?", *parentDirectoryId).Find(&newDirectory)
			oldCountDir = *directory.CountDirectories
			directory = newDirectory

		} else {
			// Дошли до домашней директории
			fmt.Println("Дошли до домашней директории")
			return nil
		}
	}

	return nil
}

func (d *Directory) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Printf("Начало работы триггера(Директория %s) BeforeDelete\n", d.Name)
	fmt.Println(d)

	// Домашняя директория
	if d.DirectoryId == nil {
		*d.CountDirectories -= 1
		return
	}

	parentDirectoryId := d.DirectoryId
	directory := &Directory{}

	tx.Where("id = ?", parentDirectoryId).First(&directory)
	// if d.DirectoryId != nil {
	// 	d.Path = ""
	// }

	directory = d

	iter := 1
	oldCountDir := 0
	for directory.Id != 0 {
		fmt.Printf("Итерация - %d\n", iter)
		iter++

		newDirectoryCount := 0
		if directory.CountDirectories != nil {
			newDirectoryCount = *directory.CountDirectories - oldCountDir - 1
		}
		directory.CountDirectories = &newDirectoryCount
		tx.Save(&directory)

		fmt.Printf("Dir name - %s. countDir - %d\n", directory.Name, *directory.CountDirectories)

		if directory.DirectoryId != nil {
			newDirectory := &Directory{}
			parentDirectoryId = directory.DirectoryId

			tx.Where("id = ?", *parentDirectoryId).Find(&newDirectory)
			oldCountDir = *directory.CountDirectories
			directory = newDirectory

		} else {
			// Дошли до домашней директории
			fmt.Println("Дошли до домашней директории")
			return nil
		}
	}

	return nil
}
