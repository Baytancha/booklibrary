package service

import (
	"fmt"
	filter "test/internal/infrastructure/filters"
	"test/internal/models"
	"testing"
)

type MockStorage struct {
}

func (m *MockStorage) Get(id int64) (*models.User, error) {
	return &models.User{}, nil
}

func (m *MockStorage) GetByName(email string) (*models.User, error) {
	return &models.User{}, nil
}

func (m *MockStorage) GetAll(filters filter.Filters) ([]*models.User, filter.Metadata, error) {
	return []*models.User{}, filter.Metadata{}, nil
}

func (m *MockStorage) Insert(user *models.User) error {
	return nil
}

func (m *MockStorage) Update(user *models.User) error {
	return nil
}

func (m *MockStorage) Delete(id int64) error {
	return nil
}

func TestUserService(t *testing.T) {
	mockStorage := MockStorage{}
	userService := NewUserService(&mockStorage)
	t.Run("ListUsers", func(t *testing.T) {
		resp, _, _ := userService.ListUsers(filter.Filters{})
		fmt.Println(resp)
	})
	t.Run("Get users by name", func(t *testing.T) {
		resp, _ := userService.GetUserByName("")
		if resp == nil {
			t.Errorf("expected error got nil")
		}

	})

	t.Run("Get users by ID", func(t *testing.T) {
		resp, _ := userService.GetUserById(0)
		if resp == nil {
			t.Errorf("expected error got nil")
		}

	})
	t.Run("Create", func(t *testing.T) {
		resp := userService.CreateUser(&models.User{})
		if resp != nil {
			t.Errorf("expected error got nil")
		}

	})

	t.Run("Update", func(t *testing.T) {
		resp := userService.UpdateUser(&models.User{})
		if resp != nil {
			t.Errorf("expected error got nil")
		}

	})

	t.Run("Delete", func(t *testing.T) {
		resp := userService.DeleteUser(0)
		if resp != nil {
			t.Errorf("expected error got nil")
		}

	})

}
