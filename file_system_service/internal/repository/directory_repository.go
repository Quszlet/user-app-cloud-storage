package repository

import (
	"fmt"
	"time"

	"github.com/Quszlet/file_system_service/internal/models"
	"gorm.io/gorm"
)

type DirectoryMySql struct {
	db *gorm.DB
}

func NewDirectoryMySql(db *gorm.DB) *DirectoryMySql {
	return &DirectoryMySql{db: db}
}

func (dmsql *DirectoryMySql) Migration() {
	dmsql.db.AutoMigrate(&models.Directory{})

	var count int64
	dmsql.db.Model(&models.Directory{}).Count(&count)

	if count == 0 {

		directory1 := models.Directory{
			DataCreate: time.Now(),
			UserId:     1,
		}

		directory2 := models.Directory{
			DataCreate: time.Now(),
			UserId:     2,
		}

		directory3 := models.Directory{
			DataCreate: time.Now(),
			UserId:     3,
			Directories: &[]models.Directory{
				{
					Name:       "Photo",
					DataCreate: time.Now(),
					UserId:     3,
				},
				{
					Name:       "Films",
					DataCreate: time.Now(),
					UserId:     3,
				},
			},
		}

		dmsql.db.Create(&directory1)
		dmsql.db.Create(&directory2)
		dmsql.db.Create(&directory3)

		idDirect5 := 5

		directory4 := &models.Directory{
			Name:        "Action",
			DataCreate:  time.Now(),
			UserId:      3,
			DirectoryId: &idDirect5,
		}

		dmsql.db.Create(&directory4)
		fmt.Println("Directories added in context migration")
	}
}

func (dmsql *DirectoryMySql) Create(directory models.Directory) (int, error) {
	directory.DataCreate = time.Now()
	result := dmsql.db.Create(&directory)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(directory.Id), nil
}

func (dmsql *DirectoryMySql) Update(directory models.Directory) error {
	dmsql.db.Debug()
	currentTime := time.Now()
	directory.DataChange = &currentTime
	fmt.Print(directory)

	parentDirId := directory.DirectoryId
	dir := &models.Directory{}
	dmsql.db.Where("id = ?", parentDirId).Find(&dir)
	directory.Path = dir.Path + "/" + directory.Name

	dmsql.db.Save(&directory)
	return nil
}

func (dmsql *DirectoryMySql) GetById(directoryId int) (models.Directory, error) {
	var directory models.Directory
	result := dmsql.db.Preload("Directories").Preload("MetadataFiles").First(&directory, directoryId)
	if result.Error != nil {
		return models.Directory{}, nil
	}
	return directory, nil
}

func (dmsql *DirectoryMySql) GetHomeDirByUserId(userId int) (models.Directory, error) {
	var directory models.Directory
	dmsql.db = dmsql.db.Debug()
	result := dmsql.db.Where("directory_id IS NULL AND user_id = ?", userId).Preload("Directories").Preload("MetadataFiles").First(&directory)
	if result.Error != nil {
		return models.Directory{}, nil
	}
	return directory, nil
}

func (dmsql *DirectoryMySql) GetByName(name string) (models.Directory, error) {
	var directory models.Directory
	result := dmsql.db.Preload("Directories").Preload("MetadataFiles").Where("name = ?", name).First(&directory)
	if result.Error != nil {
		return models.Directory{}, nil
	}
	return directory, nil
}

func (dmsql *DirectoryMySql) GetAll() ([]models.Directory, error) {
	var directories []models.Directory
	result := dmsql.db.Preload("Directories").Preload("MetadataFiles").Find(&directories)
	if result.Error != nil {
		return []models.Directory{}, nil
	}
	return directories, nil
}

func (dmsql *DirectoryMySql) Delete(directoryId int) error {
	var directory models.Directory
	dmsql.db.Select("id").First(&directory, directoryId)
	result := dmsql.db.Where("id = ?", directoryId).Delete(&directory, directoryId)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (dmsql *DirectoryMySql) Move(movedDirectoryId, directoryId int) error {
	tx1 := dmsql.db.Begin()

	fmt.Printf("movedDirectoryId - %d. directoryId - %d\n", movedDirectoryId, directoryId)

	var currentDir models.Directory
	tx1.Where("id = ?", movedDirectoryId).Find(&currentDir)

	// err := tx.Delete(&currentDir).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return fmt.Errorf("ошибка при удалении директории %v", err)
	// }

	userId := currentDir.UserId
	currentDir.DirectoryId = &directoryId

	// currentDir.Id = 0

	fmt.Println(&currentDir)
	fmt.Println(*currentDir.DirectoryId)

	// parentDir :=  &directoryId
	// if *parentDir == 0 {
	// 	parentDir
	// }

	err := tx1.Model(&models.Directory{}).Where("id = ?", movedDirectoryId).Update("directory_id", directoryId).Update("path", "").Error

	if err != nil {
		tx1.Rollback()
		return fmt.Errorf("ошибка при обновлении директории %v", err)
	}

	if err := tx1.Commit().Error; err != nil {
		return err
	}

	tx2 := dmsql.db.Begin()

	var homeDir models.Directory
	err = tx2.Where("directory_id IS NULL AND user_id = ?", userId).First(&homeDir).Error
	if err != nil {
		tx2.Rollback()
		return fmt.Errorf("ошибка при взятии домашней директории %v", err)
	}

	fmt.Println("homedir", homeDir)
	countDir := 0
	dmsql.calculateSizeDir(homeDir.Id, tx2, "", &countDir)

	if err := tx2.Commit().Error; err != nil {
		return err
	}

	fmt.Printf("Конец транзакции Move\n")
	return nil
}

func (dmsql *DirectoryMySql) calculateSizeDir(DirectoryId int, tx *gorm.DB, path string, countDir *int) int {
	fmt.Printf("Реккурсия директория - %d\n", DirectoryId)
	if DirectoryId == 0 {
		return 0
	}

	var directory models.Directory
	tx.Preload("Directories").Preload("MetadataFiles").Where("id = ?", DirectoryId).Find(&directory)

	sum := 0

	path += "/" + directory.Name
	directory.Path = path

	if directory.Directories != nil {
		fmt.Printf("calculateSizeDir: зашли в подкатологи\n ")
		directories := *directory.Directories
		fmt.Println("calculateSizeDir: все поддиректории - ")
		fmt.Println(directories)

		for _, dir := range directories {
			sum += dmsql.calculateSizeDir(dir.Id, tx, path, countDir)
		}

		*countDir += len(directories)

	}

	if directory.MetadataFiles != nil {
		files := *directory.MetadataFiles

		for _, file := range files {
			sum += file.Size
		}
	}

	directory.Size = sum
	directory.CountDirectories = countDir
	directory.Directories = nil
	directory.MetadataFiles = nil

	fmt.Println(directory)

	fmt.Printf("directory.Path - %s\n", directory.Path)
	fmt.Printf("countDir - %d", *countDir)

	tx.Save(&directory)

	return sum
}
