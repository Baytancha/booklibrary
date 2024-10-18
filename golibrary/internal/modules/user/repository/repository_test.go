package repository

import (
	filter "test/internal/infrastructure/filters"
	"test/internal/models"
	"testing"
)

func NewMockStorage() *MockStorage {
	return &MockStorage{
		GetByEmail_mock: func(email string) (*models.User, error) {
			return &models.User{}, nil
		},
		Get_mock: func(id int64) (*models.User, error) {
			return &models.User{}, nil
		},
		GetAll_mock: func(filters filter.Filters) ([]*models.User, filter.Metadata, error) {
			return []*models.User{}, filter.Metadata{}, nil
		},
		Insert_mock: func(user *models.User) error {
			return nil
		},
		Update_mock: func(user *models.User) error {
			return nil
		},
		Delete_mock: func(id int64) error {
			return nil
		},
	}
}

type MockStorage struct {
	GetByEmail_mock func(email string) (*models.User, error)
	Get_mock        func(id int64) (*models.User, error)
	GetAll_mock     func(filters filter.Filters) ([]*models.User, filter.Metadata, error)
	Insert_mock     func(user *models.User) error
	Update_mock     func(user *models.User) error
	Delete_mock     func(id int64) error
}

func (m *MockStorage) GetByEmail(email string) (*models.User, error) {
	return m.GetByEmail_mock(email)
}

func (m *MockStorage) Get(id int64) (*models.User, error) {
	return m.Get_mock(id)
}

func (m *MockStorage) GetAll(filters filter.Filters) ([]*models.User, filter.Metadata, error) {
	return m.GetAll_mock(filters)
}

func (m *MockStorage) Insert(user *models.User) error {
	return m.Insert_mock(user)
}

func (m *MockStorage) Update(user *models.User) error {
	return m.Update_mock(user)
}

func (m *MockStorage) Delete(id int64) error {
	return m.Delete_mock(id)
}

func TestRepo(t *testing.T) {
	userRepository := NewMockStorage()

	t.Run("Get user by email", func(t *testing.T) {
		resp, _ := userRepository.GetByEmail("")
		if resp == nil {
			t.Errorf("expected error got nil")
		}
	})

	t.Run("Get user by ID", func(t *testing.T) {
		resp, _ := userRepository.Get(0)
		if resp == nil {
			t.Errorf("expected error got nil")
		}
	})

	t.Run("List users", func(t *testing.T) {
		resp, _, _ := userRepository.GetAll(filter.Filters{})
		if resp == nil {
			t.Errorf("expected error got nil")
		}
	})

	t.Run("Create", func(t *testing.T) {
		resp := userRepository.Insert(&models.User{})
		if resp != nil {
			t.Errorf("expected error got nil")
		}
	})
}
