package middleware

import (
	"Deall/domain"
	"Deall/helpers"
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

type Middleware struct {
	userRepo  domain.UserRepository
	tokenRepo domain.TokenRepository
}

func NewMiddleware(u domain.UserRepository, t domain.TokenRepository) Middleware {
	handler := &Middleware{
		userRepo:  u,
		tokenRepo: t,
	}

	return *handler
}

// AuthMiddleware function to authenticate requests using the provided token
func (m *Middleware) AuthMiddleware(next httprouter.Handle, allowedRoles []string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		var (
			resp = helpers.Response{
				Status: helpers.FailMsg,
				Data:   nil,
			}
		)
		w.Header().Set("Content-Type", "application/json")

		ctx := context.Background()
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logrus.Error("Middleware | Empty auth header")
			resp.Message = "Empty Authorization header"

			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Extract the token from the header
		splitHeader := strings.Split(authHeader, " ")
		if len(splitHeader) != 2 || strings.ToLower(splitHeader[0]) != "token" {
			logrus.Error("Middleware | Empty token")
			resp.Message = "Empty token"

			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(resp)
			return
		}
		token := splitHeader[1]

		checkToken, err := m.tokenRepo.GetByToken(ctx, token)
		if err != nil {
			resp.Message = err.Error()
			switch {
			case err == mongo.ErrNoDocuments:
				resp.Message = "no user match with the given token, please check your token"

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

		usr, err := m.userRepo.GetByID(ctx, checkToken.UserID)
		if err != nil {
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

		if !domain.IsValidContains(allowedRoles, usr.Role) {
			logrus.Error("Middleware |Unauthorized roles")
			resp.Message = "Unauthorized roles"

			// Serialize the error response to JSON and send it back to the client
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(resp)
			return
		}

		// If authentication succeeded, call the next handler with the modified context
		next(w, r.WithContext(ctx), ps)
	}
}
