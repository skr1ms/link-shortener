package auth

import (
	"net/http"

	"linkshortener/configs"
	"linkshortener/pkg/req"
	"linkshortener/pkg/res"
)

type AuthHandlerDeps struct {
	Config      *configs.AuthConfig
	AuthService *AuthService
}

type AuthHandler struct {
	deps *AuthHandlerDeps
}

func NewAuthHandler(router *http.ServeMux, deps *AuthHandlerDeps) {
	authHandler := &AuthHandler{
		deps: deps,
	}
	router.HandleFunc("POST /auth/register", authHandler.Register())
	router.HandleFunc("POST /auth/login", authHandler.Login())
	router.HandleFunc("POST /auth/refresh", authHandler.RefreshToken())
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}

		user, err := handler.deps.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Response(w, 201, user)
	}
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}

		user, err := handler.deps.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		accessToken, refreshToken, err := handler.deps.AuthService.jwt.CreateTokenPair(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Response(w, 200, map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func (handler *AuthHandler) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RefreshTokenRequest](&w, r)
		if err != nil {
			return
		}

		user, err := handler.deps.AuthService.jwt.ValidateRefreshToken(body.RefreshToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		accessToken, newRefreshToken, err := handler.deps.AuthService.jwt.CreateTokenPair(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Response(w, 200, RefreshTokenResponse{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
		})
	}
}
