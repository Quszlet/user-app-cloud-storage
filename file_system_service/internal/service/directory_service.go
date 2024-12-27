package service

import (
	"github.com/Quszlet/file_system_service/internal/models"
	"github.com/Quszlet/file_system_service/internal/repository"
)

type DirectoryService struct {
	repo repository.Directory
}

func NewDirectoryService(r *repository.Repository) *DirectoryService {
	return &DirectoryService{repo: r.Directory}
}

func (ds *DirectoryService) Create(directory models.Directory) (int, error) {
	return ds.repo.Create(directory)
}

func (ds *DirectoryService) Update(directory models.Directory) error {
	return ds.repo.Update(directory)
}

func (ds *DirectoryService) GetById(directoryId int) (models.Directory, error) {
	return ds.repo.GetById(directoryId)
}

func (ds *DirectoryService) GetHomeDirByUserId(directoryId int) (models.Directory, error) {
	return ds.repo.GetHomeDirByUserId(directoryId)
}

func (ds *DirectoryService) GetByName(name string) (models.Directory, error) {
	return ds.repo.GetByName(name)
}

func (ds *DirectoryService) GetAll() ([]models.Directory, error) {
	return ds.repo.GetAll()
}

func (ds *DirectoryService) Delete(directoryId int) error {
	return ds.repo.Delete(directoryId)
}

func (ds *DirectoryService) Move(movedFileId, directoryId int) error {
	return  ds.repo.Move(movedFileId, directoryId)
}