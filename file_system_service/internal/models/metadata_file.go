package models

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"gorm.io/gorm"
)

type MetadataFile struct {
	Id          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Size        int       `json:"size" gorm:"not null"`
	DataCreate  time.Time `json:"data_create" gorm:"not null"`
	DataChange  *time.Time `json:"data_change"`
	Path        string `json:"path" gorm:"not null"`
	DirectoryId int    `json:"directory_id" gorm:"not null"` // Связь с Directory
	UserId      int    `json:"user_id" gorm:"not null"`      // Владелец файла
}

func (mf MetadataFile) Validate() error {
	return validation.ValidateStruct(&mf,
		validation.Field(&mf.Name, validation.Required),
		validation.Field(&mf.Size, validation.Required), //!!!
		validation.Field(&mf.DirectoryId, validation.Required),
		validation.Field(&mf.UserId, validation.Required))
}


func (mf *MetadataFile) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Printf("Начало работы триггера(файл %s) AfterCreate\n", mf.Name)
	fmt.Println(mf)


	err = mf.UpdateInfoAboutFileInTopDir(tx)
	if err != nil {
		fmt.Printf("Ошибка хука AfterCreate: %s", err.Error())
		return err
	}

	fmt.Println("Конец хука AfterCreate")
	return nil
}

func (mf *MetadataFile) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Printf("Начало работы триггера(файл %s) BeforeDelete\n", mf.Name)
	fmt.Println(mf)

	err = mf.ClearFileInfoInTopDir(tx)
	if err != nil {
		fmt.Printf("Ошибка хука BeforeDelete: %s", err.Error())
		return err
	}
	
	fmt.Println("Конец хука BeforeDelete")
	return nil
}

func (mf *MetadataFile) UpdateInfoAboutFileInTopDir(tx *gorm.DB) (err error) {
	currnetDirectoryId := mf.DirectoryId
	directory := &Directory{}
	iter := 1

	tx.Where("id = ?", currnetDirectoryId).First(&directory)
	mf.Path = directory.Path + "/" + mf.Name // Дополнение пути ТУТ БЫЛИ ИЗМЕНЕНИЯ
	tx.Save(&mf)
	fmt.Println(mf)

	for directory.Id != 0 {
		fmt.Printf("Итерация - %d\n", iter)
		iter++

		directory.Size += mf.Size
		countFiles := directory.CountFiles
		var newCount int

		if countFiles != nil {
			newCount = *countFiles + 1
		} else {
			newCount = 1
		}

		directory.CountFiles = &newCount

		tx.Save(&directory)
		fmt.Printf("Обновили размер и количество файлов директории с id %d\n", directory.Id)
		fmt.Println(directory)

		if directory.DirectoryId != nil {
			newDirectory := &Directory{}
			parentDirectoryId := *directory.DirectoryId

			tx.Where("id = ?", parentDirectoryId).Find(&newDirectory)
			directory = newDirectory
			fmt.Printf("Перешли к родительской директории с id %d\n", directory.Id)
		} else {
			// Дошли до домашней директории
			fmt.Println("Дошли до домашней директории")
			return nil
		}
	}

	return nil
}

func (mf *MetadataFile)  ClearFileInfoInTopDir(tx *gorm.DB) (err error) {
	currnetDirectoryId := mf.DirectoryId
	directory := &Directory{}
	iter := 1
	fmt.Println(mf)
	fmt.Printf("currnetDirectoryId - %d", currnetDirectoryId)
	tx.Where("id = ?", currnetDirectoryId).First(&directory)
	mf.Path = ""
	fmt.Println(mf)

	// Поднимаемся и подчищаем поля, связанные с файлом
	for directory.Id != 0 {
		fmt.Printf("Итерация - %d\n", iter)
		iter++

		directory.Size -= mf.Size
		countFiles := directory.CountFiles
		var newCount int

		if countFiles != nil {
			newCount = *countFiles - 1
		} else {
			newCount = 0
		}
		directory.CountFiles = &newCount

		tx.Save(&directory)
		fmt.Printf("Уменьшили размер и количество файлов директории с id %d\n", directory.Id)
		fmt.Println(directory)

		if directory.DirectoryId != nil {
			newDirectory := &Directory{}
			parentDirectoryId := *directory.DirectoryId

			tx.Where("id = ?", parentDirectoryId).Find(&newDirectory)
			directory = newDirectory
			fmt.Printf("Перешли к родительской директории с id %d\n", directory.Id)
		} else {
			// Дошли до домашней директории
			fmt.Println("Дошли до домашней директории")
			return nil
		}
	}

	return nil
}
