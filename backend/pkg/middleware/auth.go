package middleware

import (
	"context"
	"net/http"
	"strings"

	"linkshortener/config"
	"linkshortener/pkg/jwt"
)

func writeUnauthorized(w http.ResponseWriter) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

func IsAuthenticated(next http.Handler, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			writeUnauthorized(w)
			return
		}
		if !strings.HasPrefix(token, "Bearer ") {
			writeUnauthorized(w)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		jwtService := jwt.NewJWT(config.Auth.SecretKey, config.Auth.RefreshTokenSecretKey)

		user, err := jwtService.ValidateToken(token)
		if err != nil {
			refreshToken := r.Header.Get("X-Refresh-Token")
			if refreshToken == "" {
				writeUnauthorized(w)
				return
			}

			newUser, newAccessToken, newRefreshToken, refreshErr := jwtService.RefreshTokens(refreshToken)
			if refreshErr != nil {
				writeUnauthorized(w)
				return
			}

			w.Header().Set("X-New-Access-Token", newAccessToken)
			w.Header().Set("X-New-Refresh-Token", newRefreshToken)

			ctx := context.WithValue(r.Context(), "userEmail", newUser.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		ctx := context.WithValue(r.Context(), "userEmail", user.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
