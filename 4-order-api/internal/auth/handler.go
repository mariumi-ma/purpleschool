package auth

import (
	"errors"
	"net/http"

	"purpleschool/configs"
	"purpleschool/internal/request"
	"purpleschool/internal/response"
)

type AuthHandler struct {
	Config  *configs.Config
	Service *AuthService
}

func NewAuthHandler(router *http.ServeMux, config *configs.Config, service *AuthService) {
	handler := &AuthHandler{
		Config:  config,
		Service: service,
	}

	router.Handle("POST /auth", handler.Auth())
	router.Handle("POST /auth/verify", handler.VerifyCode())
}

func (h *AuthHandler) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[UserAuthRequest](w, r)
		if err != nil {
			return
		}

		sessionID, err := h.Service.Auth(body.Phone)
		if err != nil {
			response.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := UserAuthResponse{
			SessionID: sessionID,
		}

		response.JSON(w, resp, http.StatusOK)
	}
}

func (h *AuthHandler) VerifyCode() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := request.HandleBody[UserAuthRequest](w, r)
		if err != nil {
			return
		}

		user, err := h.Service.VerifyCode(body.SessionID, body.Code)
		if err != nil {
			if errors.Is(err, ErrWrongCredentials) {
				response.JSON(w, err.Error(), http.StatusUnauthorized)
				return
			}
			response.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := NewJWT("secret").GenerateToken(user.SessionID)
		if err != nil {
			response.JSON(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := Token{
			Token: token,
		}

		response.JSON(w, resp, http.StatusOK)
	}
}
