package service

import (
	"errors"
	"test/internal/infrastructure/filters"
	"test/internal/infrastructure/validator"
	"test/internal/modules/user/repository"

	"test/internal/models"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrNoUser         = errors.New("user doesn't exist")
	ErrWrongPassword  = errors.New("wrong password")
)

// r.Post("/user", ctrl.UserHandler.CreateUser)//ctrl.Auth.Register
// 	r.Get("/user/login", ctrl.Auth.Login)
// 	r.Get("/user/logout", ctrl.Auth.Logout)
// 	r.Get("/user/{username}", ctrl.UserHandler.GetUserByEmail)
// 	r.Put("/user/{username}", ctrl.UserHandler.UpdateUser)
// 	r.Delete("/user/{username}", ctrl.UserHandler.DeleteUser)
// 	r.Post("/user/CreateWithList", ctrl.UserHandler.ListUsers)
// 	r.Post("/user/CreateWithArray", ctrl.UserHandler.CreateArray)

type IUserService interface {
	GetUserByName(email string) (*models.User, error)
	GetUserById(id int64) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(id int64) error
	ListUsers(filters filters.Filters) ([]*models.User, filters.Metadata, error)
}

type UserService struct {
	storage repository.IUserStorage
}

func NewUserService(repo repository.IUserStorage) *UserService {
	return &UserService{storage: repo}
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.storage.Insert(user)
}

func (s *UserService) GetUserByName(username string) (*models.User, error) {
	return s.storage.GetByName(username)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.storage.Update(user)
}

func (s *UserService) GetUserById(id int64) (*models.User, error) {
	return s.storage.Get(id)
}

func (s *UserService) DeleteUser(id int64) error {
	return s.storage.Delete(id)
}

func (s *UserService) ListUsers(filters filters.Filters) ([]*models.User, filters.Metadata, error) {
	return s.storage.GetAll(filters)
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *models.User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.Plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.Plaintext)
	}

	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}
