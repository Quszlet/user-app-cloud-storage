package repository

import (
	"fmt"
	"time"

	"github.com/Quszlet/file_system_service/internal/models"
	"gorm.io/gorm"
)

type MetadataFileMySql struct {
	db *gorm.DB
}

func NewMetadataFileMySql(db *gorm.DB) *MetadataFileMySql {
	return &MetadataFileMySql{db: db}
}

func (mdfmsql *MetadataFileMySql) Migration() {
	mdfmsql.db.AutoMigrate(&models.MetadataFile{})

	var count int64
	mdfmsql.db.Model(&models.MetadataFile{}).Count(&count)

	if count == 0 {
		metadataFiles := []models.MetadataFile{
			{
				Name:        "ADs.txt",
				Size:        1024,
				DataCreate:  time.Now(),
				DirectoryId: 1,
				UserId:      1,
			},
			{
				Name:        "Check.pdf",
				Size:        10,
				DataCreate:  time.Now(),
				DirectoryId: 2,
				UserId:      2,
			},
			{
				Name:        "root3",
				Size:        64,
				DataCreate:  time.Now(),
				UserId:      3,
				DirectoryId: 4,
			},
			{
				Name:        "fil",
				Size:        128,
				DataCreate:  time.Now(),
				UserId:      3,
				DirectoryId: 4,
			},
			{
				Name:        "чай.txt",
				Size:        1024,
				DataCreate:  time.Now(),
				UserId:      3,
				DirectoryId: 5,
			},
		}
		mdfmsql.db.Create(&metadataFiles)
		fmt.Println("Directories added in context migration")
	}
}

func (mdfmsql *MetadataFileMySql) Create(metadataFile models.MetadataFile) (int, error) {
	result := mdfmsql.db.Create(&metadataFile)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(metadataFile.Id), nil
}

func (mdfmsql *MetadataFileMySql) Update(metadataFile models.MetadataFile) error {
	currentTime := time.Now()
	metadataFile.DataChange = &currentTime
	fmt.Print(metadataFile)

	parentDirId := metadataFile.DirectoryId
	dir := &models.Directory{}
	mdfmsql.db.Where("id = ?", parentDirId).Find(&dir)
	metadataFile.Path = dir.Path + "/" + metadataFile.Name

	mdfmsql.db.Save(&metadataFile)
	return nil
}

func (mdfmsql *MetadataFileMySql) GetById(metadataFileId int) (models.MetadataFile, error) {
	var metadataFile models.MetadataFile
	result := mdfmsql.db.First(&metadataFile, metadataFileId)
	if result.Error != nil {
		return models.MetadataFile{}, nil
	}
	return metadataFile, nil
}

func (mdfmsql *MetadataFileMySql) GetByName(name string) (models.MetadataFile, error) {
	var metadataFile models.MetadataFile
	result := mdfmsql.db.Where("name = ?", name).First(&metadataFile)
	if result.Error != nil {
		return models.MetadataFile{}, nil
	}
	return metadataFile, nil
}

func (mdfmsql *MetadataFileMySql) GetAll() ([]models.MetadataFile, error) {
	var directories []models.MetadataFile
	result := mdfmsql.db.Find(&directories)
	if result.Error != nil {
		return []models.MetadataFile{}, nil
	}
	return directories, nil
}

func (mdfmsql *MetadataFileMySql) Delete(metadataFileId int) error {
	var file models.MetadataFile
	mdfmsql.db.Select("directory_id, size").First(&file, metadataFileId)
	result := mdfmsql.db.Where("id = ?", metadataFileId).Delete(&file)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (mdfmsql *MetadataFileMySql) Move(movedFileId, directoryId int) error {
	tx := mdfmsql.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var file models.MetadataFile
	tx.First(&file, movedFileId)

	if err := tx.Delete(&file).Error; err != nil {
		tx.Rollback()
		return err
	}

	file.DirectoryId = directoryId
	file.Path = ""

	if directoryId == 0 {
		homeDir := models.Directory{}
		tx.Where("directory_id IS NULL AND user_id = ?", file.UserId).First(&homeDir)

		file.DirectoryId = homeDir.Id
	}

	if err := tx.Create(&file).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
