package handler

import (
	"Deall/domain"
	"Deall/helpers"
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
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
			Username: r.PostFormValue("username"),
			Password: r.PostFormValue("password"),
			Role:     r.PostFormValue("role"),
		}
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

	users, err = u.userService.GetUsers(ctx)
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

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var (
		err  error
		resp = helpers.Response{
			Status:  helpers.SuccessMsg,
			Message: "",
			Data:    nil,
		}
		user domain.Users
		usr  = r.URL.Query().Get("username")
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	if usr == "" {
		resp.Status = helpers.FailMsg
		resp.Message = "username cannot be empty"

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	user, err = u.userService.GetUser(ctx, usr)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Message = err.Error()
		switch {
		case err == mongo.ErrNoDocuments:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		default:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	resp.Data = user
	json.NewEncoder(w).Encode(resp)
	return
}

func (u *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		err  error
		resp = helpers.Response{
			Status:  helpers.SuccessMsg,
			Message: "",
			Data:    nil,
		}
		id = ps.ByName("id")
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	if id == "" {
		resp.Status = helpers.FailMsg
		resp.Message = "id cannot be empty"

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = u.userService.DeleteUser(ctx, id)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Message = err.Error()
		switch {
		case err == mongo.ErrNoDocuments:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		default:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	json.NewEncoder(w).Encode(resp)
	return
}

func (u *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		err  error
		resp = helpers.Response{
			Status:  helpers.SuccessMsg,
			Message: "Data have been updated!",
			Data:    nil,
		}
		id   = ps.ByName("id")
		user = domain.Users{
			Username: r.PostFormValue("username"),
			Password: r.PostFormValue("password"),
			Role:     r.PostFormValue("role"),
		}
	)
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	if id == "" {
		resp.Status = helpers.FailMsg
		resp.Message = "id cannot be empty"

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if user.Role != "" {
		if err = user.ValidateRole(); err != nil {
			resp.Status = helpers.FailMsg
			resp.Message = err.Error()

			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	user.ID = id
	err = u.userService.UpdateUser(ctx, user)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Message = err.Error()
		switch {
		case err == mongo.ErrNoDocuments:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		default:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	json.NewEncoder(w).Encode(resp)
	return
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var (
		err  error
		resp = helpers.Response{
			Status:  helpers.SuccessMsg,
			Message: "Here's your login token:",
			Data:    nil,
		}
		username = r.PostFormValue("username")
		pwd      = r.PostFormValue("password")
		token    string
	)
	w.Header().Set("Content-Type", "application/json")

	if username == "" || pwd == "" {
		resp.Status = helpers.FailMsg
		resp.Message = "username or password cannot be empty"

		// Serialize the error response to JSON and send it back to the client
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), u.timeout)
	defer cancel()

	token, err = u.userService.Authentication(ctx, username, pwd)
	if err != nil {
		resp.Status = helpers.FailMsg
		resp.Message = err.Error()
		switch {
		case err == mongo.ErrNoDocuments:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		default:
			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	resp.Data = map[string]interface{}{"token": token}
	json.NewEncoder(w).Encode(resp)
	return
}
