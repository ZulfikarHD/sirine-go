package services

import (
	"sirine-go/backend/database"
	"sirine-go/backend/models"
)

type ExampleService struct{}

func NewExampleService() *ExampleService {
	return &ExampleService{}
}

func (s *ExampleService) GetAll() ([]models.Example, error) {
	var examples []models.Example
	result := database.GetDB().Find(&examples)
	return examples, result.Error
}

func (s *ExampleService) GetByID(id uint) (*models.Example, error) {
	var example models.Example
	result := database.GetDB().First(&example, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &example, nil
}

func (s *ExampleService) Create(example *models.Example) error {
	return database.GetDB().Create(example).Error
}

func (s *ExampleService) Update(id uint, example *models.Example) error {
	return database.GetDB().Model(&models.Example{}).Where("id = ?", id).Updates(example).Error
}

func (s *ExampleService) Delete(id uint) error {
	return database.GetDB().Delete(&models.Example{}, id).Error
}
