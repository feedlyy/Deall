package handler

import (
	"Deall/domain"
	"Deall/helpers"
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

type UserHandler struct {
	userService domain.UserService
	timeout     time.Duration
}

func NewUserHandler(u domain.UserService, timeout time.Duration) UserHandler {
	handler := &UserHandler{
		userService: u,
		timeout:     timeout,
	}

	return *handler
}

func (u *UserHandler) RegistUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err  error
		resp = helpers.Response{
			Status:  helpers.SuccessMsg,
			Message: "user have been created!",
			Data:    nil,
		}
		user = domain.Users{
			Username:  r.PostFormValue("username"),
			Password:  r.PostFormValue("password"),
			Role:      r.PostFormValue("role"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		pwd string
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	if err = user.ValidateUser(); err != nil {
		resp.Status = helpers.FailMsg
		resp.Message = err.Error()

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// bcrypt the password
	pwd, err = domain.HashPassword(user.Password)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Message = err.Error()

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	user.Password = pwd

	err = u.userService.Register(ctx, user)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Message = err.Error()

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	json.NewEncoder(w).Encode(resp)
	return
}

func (u *UserHandler) GetAllUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err  error
		resp = helpers.Response{
			Status:  helpers.SuccessMsg,
			Message: "",
			Data:    nil,
		}
		users []domain.Users
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	users, err = u.userService.GetUSers(ctx)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Message = err.Error()

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Data = users
	json.NewEncoder(w).Encode(resp)
	return
}
