package service

import (
	"github.com/Quszlet/file_system_service/internal/models"
	"github.com/Quszlet/file_system_service/internal/repository"
)

type Directory interface {
	Create(directory models.Directory) (int, error)
	Update(directory models.Directory) error
	GetById(diretoryId int) (models.Directory, error)
	GetByName(name string) (models.Directory, error)
	GetHomeDirByUserId(userId int) (models.Directory, error)
	GetAll() ([]models.Directory, error)
	Delete(diretoryId int) error
	Move(movedDirectoryId, directoryId int) error
}

type MetadataFile interface {
	Create(metadataFile models.MetadataFile) (int, error)
	Update(metadataFile models.MetadataFile) error
	GetById(metadataFileId int) (models.MetadataFile, error)
	GetByName(name string) (models.MetadataFile, error)
	GetAll() ([]models.MetadataFile, error)
	Delete(metadataFileId int) error
	Move(movedFileId, directoryId int) error
}

type Service struct {
	Directory
	MetadataFile
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		Directory: NewDirectoryService(r),
		MetadataFile: NewMetadataFileService(r),
	}
}