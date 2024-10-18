package repository

import (
	filter "test/internal/infrastructure/filters"
	"test/internal/models"
)

type IUserStorage interface {
	GetByName(username string) (*models.User, error)
	Get(id int64) (*models.User, error)
	GetAll(filters filter.Filters) ([]*models.User, filter.Metadata, error)
	Insert(user *models.User) error
	Update(user *models.User) error
	Delete(id int64) error
}
