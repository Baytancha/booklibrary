package controller

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	filter "test/internal/infrastructure/filters"
	"test/internal/infrastructure/responder"
	"test/internal/models"
	"test/internal/modules/user/service"
	"testing"

	"github.com/go-chi/chi"
	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

type MockStorage struct {
	GetByName_mock func(email string) (*models.User, error)
	Get_mock       func(id int64) (*models.User, error)
	GetAll_mock    func(filters filter.Filters) ([]*models.User, filter.Metadata, error)
	Insert_mock    func(user *models.User) error
	Update_mock    func(user *models.User) error
	Delete_mock    func(id int64) error
}

func (m *MockStorage) GetByName(email string) (*models.User, error) {
	return m.GetByName_mock(email)
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

// r.Post("/user", ctrl.UserHandler.CreateUser)
// 	r.Get("/user/login", ctrl.UserHandler.Login)
// 	r.Get("/user/logout", ctrl.UserHandler.Logout)
// 	r.Get("/user/{username}", ctrl.UserHandler.GetUserByName)
// 	r.Put("/user/{username}", ctrl.UserHandler.UpdateUser)
// 	r.Delete("/user/{username}", ctrl.UserHandler.DeleteUser)
// 	r.Post("/user/CreateWithList", ctrl.UserHandler.CreateWithList)
// 	r.Post("/user/CreateWithArray", ctrl.UserHandler.CreateWithArray)

func TestCreateHandler(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/user", bytes.NewReader([]byte(`{ "name": "bigsmoke","email": "trest@example.com","password": "123456789"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Insert_mock: func(user *models.User) error { return nil },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.CreateUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})
	t.Run("dublicate", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/user", bytes.NewReader([]byte(`{ "name": "bigsmoke","email": "trest@example.com","password": "123456789"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Insert_mock: func(user *models.User) error { return service.ErrDuplicateEmail },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.CreateUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

	t.Run("internal server error", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/user", bytes.NewReader([]byte(`{ "name": "bigsmoke","email": "trest@example.com","password": "123456789"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Insert_mock: func(user *models.User) error { return errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.CreateUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

	t.Run("invalid json", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/user", bytes.NewReader([]byte(`{ "nae:,, "bigsmoke","eil": "tree.com","password": "123456789"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Insert_mock: func(user *models.User) error { return errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.CreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

	t.Run("invalid password request", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/user", bytes.NewReader([]byte(`{ "name": "bigsmoke","email": "treee.com","password": "123456789"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Insert_mock: func(user *models.User) error { return errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.CreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

	t.Run("invalid password request", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/user", bytes.NewReader([]byte(`{ "name": "bigsmoke","email": "treee.com","password": "123456789"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Insert_mock: func(user *models.User) error { return errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.CreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

}

func TestGetByIDHandler(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/user/1", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Get_mock: func(id int64) (*models.User, error) { return &models.User{}, nil },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("id", "1")
		controller.GetUserById(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})
	t.Run("user not found", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/user/", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Get_mock: func(id int64) (*models.User, error) { return nil, service.ErrRecordNotFound },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("id", "1")
		controller.GetUserById(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

	t.Run("internal server error", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/user/1", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Get_mock: func(id int64) (*models.User, error) { return nil, errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("id", "1")
		controller.GetUserById(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

	t.Run("invalid request", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/user/lklj", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Get_mock: func(id int64) (*models.User, error) { return nil, errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("id", "lklj")
		controller.GetUserById(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

}

func TestGetByEmailHandler(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/user", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetByName_mock: func(email string) (*models.User, error) { return &models.User{}, nil },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("username", "ljlj")

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)

		controller.GetUserByName(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})
	t.Run("invalid request", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/user/", nil)
		req.Header.Set("Content-Type", "application/json")
		req.URL.RawQuery = "emil=trest@example.com"

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Get_mock: func(id int64) (*models.User, error) { return nil, service.ErrRecordNotFound },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.GetUserById(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

	t.Run("invalid email", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/user/", nil)
		req.Header.Set("Content-Type", "application/json")
		req.URL.RawQuery = "eml=tres@xample.com"

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Get_mock: func(id int64) (*models.User, error) { return nil, service.ErrRecordNotFound },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.GetUserById(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

	t.Run("internal error", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/user/", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetByName_mock: func(email string) (*models.User, error) { return nil, service.ErrDuplicateEmail },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("username", "ljlj")

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)

		controller.GetUserByName(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

}

func TestUpdateHandler(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("PUT", "/user/1", bytes.NewReader([]byte(`{ "name": "bigsmoke","email": "trest@example.com","password": "123456789"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock_user := &models.User{
			ID:    1,
			Name:  "bigsmoke",
			Email: "trest@example.com",
			Password: models.Password{
				Hash: []byte("123456789"),
			},
		}

		mock := &MockStorage{
			GetByName_mock: func(email string) (*models.User, error) { return mock_user, nil },
			Update_mock:    func(user *models.User) error { return nil },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("username", "john")
		controller.UpdateUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})
	t.Run("dublicate", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/user/1", bytes.NewReader([]byte(`{ "name": "bigsmoke","email": "trest@example.com","password": "123456789"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{

			Insert_mock: func(user *models.User) error { return service.ErrDuplicateEmail },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.CreateUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

	t.Run("internal server error", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/api/users", bytes.NewReader([]byte(`{ "name": "bigsmoke","email": "trest@example.com","password": "123456789"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Insert_mock: func(user *models.User) error { return errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.CreateUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

	t.Run("invalid json", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/api/users", bytes.NewReader([]byte(`{ "nae:,, "bigsmoke","eil": "tree.com","password": "123456789"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Insert_mock: func(user *models.User) error { return errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.CreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

	t.Run("invalid password request", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/api/users", bytes.NewReader([]byte(`{ "name": "bigsmoke","email": "treee.com","password": "123456789"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Insert_mock: func(user *models.User) error { return errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.CreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

	t.Run("invalid password request", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/api/users", bytes.NewReader([]byte(`{ "name": "bigsmoke","email": "treee.com","password": "123456789"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			Insert_mock: func(user *models.User) error { return errors.New("some error") },
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.CreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

}

func TestListAllHandler(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/user/list", nil)
		req.Header.Set("Content-Type", "application/json")
		req.URL.RawQuery = "page_size=3&page=3"

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetAll_mock: func(filters filter.Filters) ([]*models.User, filter.Metadata, error) {
				return []*models.User{}, filter.Metadata{}, nil
			},
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.ListUsers(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})
	t.Run("internal server error", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/api/users/list", nil)
		req.Header.Set("Content-Type", "application/json")
		req.URL.RawQuery = "page_size=999999"

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetAll_mock: func(filters filter.Filters) ([]*models.User, filter.Metadata, error) {
				return []*models.User{}, filter.Metadata{}, nil
			},
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.ListUsers(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

	t.Run("internal server error", func(t *testing.T) {

		req := httptest.NewRequest("POST", "/api/users", bytes.NewReader([]byte(`{ "name": "bigsmoke","email": "trest@example.com","password": "123456789"}`)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetAll_mock: func(filters filter.Filters) ([]*models.User, filter.Metadata, error) {
				return nil, filter.Metadata{}, errors.New("some error")
			},
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		controller.ListUsers(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

}

func TestDeleteHandler(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {

		req := httptest.NewRequest("DELETE", "/user", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetByName_mock: func(email string) (*models.User, error) {
				return &models.User{}, nil
			},
			Delete_mock: func(id int64) error {
				return nil
			},
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("username", "user")
		controller.DeleteUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status code %d but got %d", http.StatusOK, w.Code)
		}

	})
	t.Run("invalid request body", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/api/users/", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetByName_mock: func(email string) (*models.User, error) {
				return &models.User{}, nil
			},
			Delete_mock: func(id int64) error {
				return nil
			},
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("none", "d")
		controller.DeleteUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d but got %d", http.StatusBadRequest, w.Code)
		}

	})

	t.Run("internal server error", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/api/users/1", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetByName_mock: func(email string) (*models.User, error) {
				return &models.User{}, nil
			},
			Delete_mock: func(id int64) error {
				return errors.New("some error")
			},
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("username", "user")
		controller.DeleteUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

	t.Run("dublicate request", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/api/users/lklj", nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		mock := &MockStorage{
			GetByName_mock: func(email string) (*models.User, error) {
				return &models.User{}, nil
			},
			Delete_mock: func(id int64) error {
				return service.ErrDuplicateEmail
			},
		}

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}
		decoder := godecoder.NewDecoder(jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			ValidateJsonRawMessage: true,
			DisallowUnknownFields:  true,
		})

		service := service.NewUserService(mock)

		controller := NewUserHandler(responder.NewResponder(decoder, logger), service)
		chiCtx := chi.NewRouteContext()

		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

		chiCtx.URLParams.Add("username", "user")
		controller.DeleteUser(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, w.Code)
		}

	})

}
