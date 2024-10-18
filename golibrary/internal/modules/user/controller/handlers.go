package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"test/internal/infrastructure/filters"
	"test/internal/infrastructure/helpers"
	"test/internal/infrastructure/responder"
	"test/internal/infrastructure/validator"
	"test/internal/models"
	service "test/internal/modules/user/service"
	"time"

	"github.com/go-chi/chi"
)

// r.Post("/user", ctrl.UserHandler.CreateUser)//ctrl.Auth.Register
// 	r.Get("/user/login", ctrl.Auth.Login)
// 	r.Get("/user/logout", ctrl.Auth.Logout)
// 	r.Get("/user/{username}", ctrl.UserHandler.GetUserByEmail)
// 	r.Put("/user/{username}", ctrl.UserHandler.UpdateUser)
// 	r.Delete("/user/{username}", ctrl.UserHandler.DeleteUser)
// 	r.Post("/user/CreateWithList", ctrl.UserHandler.ListUsers)
// 	r.Post("/user/CreateWithArray", ctrl.UserHandler.CreateArray)

type IUserHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	GetUserByName(w http.ResponseWriter, r *http.Request)
	GetUserById(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	CreateWithList(w http.ResponseWriter, r *http.Request)
	CreateWithArray(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type UserHandler struct {
	responder responder.Responder
	service   service.IUserService
}

func NewUserHandler(responder responder.Responder, service service.IUserService) *UserHandler {
	return &UserHandler{
		responder: responder,
		service:   service,
	}
}

func (uc *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		uc.responder.ErrorBadRequest(w, errors.New("failed to parse body"))
		return
	}
	userName := r.Form.Get("username")
	userPassword := r.Form.Get("password")
	fmt.Println("data", userName, userPassword)

	if userName == "" || userPassword == "" {
		http.Error(w, "Missing username or password.", http.StatusBadRequest)
		return
	}
	user, err := uc.service.GetUserByName(userName)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			uc.responder.ErrorInternal(w, errors.New("User not found"))
		default:
			uc.responder.ErrorInternal(w, errors.New("Internal server error"))
		}
		return
	}

	ok, err := user.Password.Matches(userPassword)
	if !ok && err != nil {
		uc.responder.ErrorBadRequest(w, err)
		return
	} else if !ok && err == nil {
		uc.responder.ErrorBadRequest(w, errors.New("wrong password"))
		return
	}
	token := helpers.GenerateToken(userName)
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(60 * time.Minute),
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt",
		Path:     "/",
		Value:    token,
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, fmt.Sprint("successfully logged in"))

}

func (uc *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt",
		Path:     "/",
		Value:    "",
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, fmt.Sprint("successfully logged out"))
}

func (uc *UserHandler) GetUserByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "username")
	if name == "" {
		uc.responder.ErrorBadRequest(w, errors.New("Invalid request body"))
		return
	}

	fmt.Println(name)
	user, err := uc.service.GetUserByName(name)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			uc.responder.ErrorInternal(w, errors.New("User not found"))
		default:
			uc.responder.ErrorInternal(w, errors.New("Internal server error"))
		}
		return
	}

	uc.responder.OutputJSON(w, user)
}

func (uc *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		uc.responder.ErrorBadRequest(w, errors.New("Invalid request body"))
		return
	}
	id64, err := strconv.Atoi(id)
	if err != nil {
		uc.responder.ErrorBadRequest(w, errors.New("Invalid request body"))
		return
	}

	user, err := uc.service.GetUserById(int64(id64))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			uc.responder.ErrorInternal(w, errors.New("User not found"))
		default:
			uc.responder.ErrorInternal(w, errors.New("Internal server error"))
		}
		return
	}

	uc.responder.OutputJSON(w, user)
}

func (uc *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		uc.responder.ErrorBadRequest(w, errors.New("Invalid request body1"))
		return
	}
	user := &models.User{
		Name:    input.Name,
		Email:   input.Email,
		Deleted: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		uc.responder.ErrorBadRequest(w, errors.New("Invalid request body2"))
		return
	}

	v := validator.New()

	if service.ValidateUser(v, user); !v.Valid() {
		uc.responder.ErrorBadRequest(w, errors.New("Invalid request body3"))
		return
	}

	err = uc.service.CreateUser(user)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrDuplicateEmail):
			uc.responder.ErrorInternal(w, errors.New("User already exists"))
		default:
			uc.responder.ErrorInternal(w, errors.New("Internal server error"))
		}
		return
	}

	uc.responder.OutputJSON(w, user)
}

func (uc *UserHandler) CreateWithList(w http.ResponseWriter, r *http.Request) {
	var users []*models.User
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for _, user := range users {

		var input struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			uc.responder.ErrorBadRequest(w, errors.New("Invalid request body1"))
			return
		}
		user = &models.User{
			Name:    input.Name,
			Email:   input.Email,
			Deleted: false,
		}

		err = user.Password.Set(input.Password)
		if err != nil {
			uc.responder.ErrorBadRequest(w, errors.New("Invalid request body2"))
			return
		}

		v := validator.New()

		if service.ValidateUser(v, user); !v.Valid() {
			uc.responder.ErrorBadRequest(w, errors.New("Invalid request body3"))
			return
		}

		err = uc.service.CreateUser(user)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrDuplicateEmail):
				uc.responder.ErrorInternal(w, errors.New("User already exists"))
			default:
				uc.responder.ErrorInternal(w, errors.New("Internal server error"))
			}
			return
		}
	}

	uc.responder.OutputJSON(w, users)
}

func (uc *UserHandler) CreateWithArray(w http.ResponseWriter, r *http.Request) {
	var users []*models.User
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	for _, user := range users {

		var input struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			uc.responder.ErrorBadRequest(w, errors.New("Invalid request body1"))
			return
		}
		user = &models.User{
			Name:  input.Name,
			Email: input.Email,

			Deleted: false,
		}

		err = user.Password.Set(input.Password)
		if err != nil {
			uc.responder.ErrorBadRequest(w, errors.New("Invalid request body2"))
			return
		}

		v := validator.New()

		if service.ValidateUser(v, user); !v.Valid() {
			uc.responder.ErrorBadRequest(w, errors.New("Invalid request body3"))
			return
		}

		err = uc.service.CreateUser(user)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrDuplicateEmail):
				uc.responder.ErrorInternal(w, errors.New("User already exists"))
			default:
				uc.responder.ErrorInternal(w, errors.New("Internal server error"))
			}
			return
		}
	}

	uc.responder.OutputJSON(w, users)
}

func (uc *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "username")
	if name == "" {
		uc.responder.ErrorBadRequest(w, errors.New("Invalid request body1"))
		return
	}

	user, err := uc.service.GetUserByName(name)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			uc.responder.ErrorInternal(w, errors.New("User not found"))
		default:
			uc.responder.ErrorInternal(w, errors.New("Internal server error"))
		}
		return
	}
	var input struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		uc.responder.ErrorBadRequest(w, errors.New("Invalid request body3"))
		return
	}
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = *input.Email
	}

	v := validator.New()

	if service.ValidateUser(v, user); !v.Valid() {
		uc.responder.ErrorInternal(w, errors.New("Invalid request body4"))
		return
	}

	err = uc.service.UpdateUser(user)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEditConflict):
			uc.responder.ErrorInternal(w, errors.New("User not found"))
		default:
			uc.responder.ErrorInternal(w, errors.New("Internal server error"))
		}
		return
	}

	uc.responder.OutputJSON(w, user)
}

func (uc *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user := chi.URLParam(r, "username")
	if user == "" {
		uc.responder.ErrorBadRequest(w, errors.New("Invalid request body"))
		return
	}
	deleted, err := uc.service.GetUserByName(user)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			uc.responder.ErrorInternal(w, errors.New("User not found"))
		default:
			uc.responder.ErrorInternal(w, errors.New("Internal server error"))
		}
		return
	}
	err = uc.service.DeleteUser(deleted.ID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRecordNotFound):
			uc.responder.ErrorInternal(w, errors.New("User not found"))
		default:
			uc.responder.ErrorInternal(w, errors.New("Internal server error"))
		}
		return
	}

	uc.responder.OutputJSON(w, "User deleted successfully")
}

func (uc *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string
		Email string
		filters.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Filters.Page = helpers.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = helpers.ReadInt(qs, "page_size", 20, v)

	input.Filters.Sort = helpers.ReadString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "email", "-id", "-name", "-email"}

	if filters.ValidateFilters(v, input.Filters); !v.Valid() {
		uc.responder.ErrorInternal(w, errors.New("Internal server error1"))
		return
	}

	users, metadata, err := uc.service.ListUsers(input.Filters)
	if err != nil {
		uc.responder.ErrorInternal(w, errors.New("Internal server error2"))
		return
	}
	uc.responder.OutputJSON(w, map[string]interface{}{"metadata": metadata, "data": users})
}
