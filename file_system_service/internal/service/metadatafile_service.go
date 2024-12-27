package service

import (
	"github.com/Quszlet/file_system_service/internal/models"
	"github.com/Quszlet/file_system_service/internal/repository"
)

type MetadataFileService struct {
	repo repository.MetadataFile
}

func NewMetadataFileService(r *repository.Repository) *MetadataFileService {
	return &MetadataFileService{repo: r.MetadataFile}
}

func (mfs *MetadataFileService) Create(metadataFile models.MetadataFile) (int, error) {
	return mfs.repo.Create(metadataFile)
}

func (mfs *MetadataFileService) Update(metadataFile models.MetadataFile) error {
	return mfs.repo.Update(metadataFile)
}

func (mfs *MetadataFileService) GetById(metadataFileId int) (models.MetadataFile, error) {
	return mfs.repo.GetById(metadataFileId)
}

func (mfs *MetadataFileService) GetByName(name string) (models.MetadataFile, error) {
	return mfs.repo.GetByName(name)
}

func (mfs *MetadataFileService) GetAll() ([]models.MetadataFile, error) {
	return mfs.repo.GetAll()
}

func (mfs *MetadataFileService) Delete(metadataFileId int) error {
	return mfs.repo.Delete(metadataFileId)
}

func (mfs *MetadataFileService) Move(movedFileId, directoryId int) error {
	return  mfs.repo.Move(movedFileId, directoryId)
}