package jwt

import (
	"errors"
	"time"

	"linkshortener/internal/user"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	SecretKey             string
	RefreshTokenSecretKey string
}

func NewJWT(secretKey, refreshTokenSecretKey string) *JWT {
	return &JWT{
		SecretKey:             secretKey,
		RefreshTokenSecretKey: refreshTokenSecretKey,
	}
}

// Создание access токена (короткий срок жизни)
func (j *JWT) CreateToken(user *user.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.Email,
		"type":    "access",
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// Создание refresh токена (длительный срок жизни)
func (j *JWT) CreateRefreshToken(user *user.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.Email,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(j.RefreshTokenSecretKey))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

// Создание пары токенов
func (j *JWT) CreateTokenPair(user *user.User) (string, string, error) {
	accessToken, err := j.CreateToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := j.CreateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Валидация access токена
func (j *JWT) ValidateToken(accessToken string) (*user.User, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.SecretKey), nil
	})

	user := &user.User{}
	if err != nil {
		return user, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if tokenType, ok := claims["type"].(string); !ok || tokenType != "access" {
			return user, errors.New("invalid token type")
		}

		if email, ok := claims["user_id"].(string); ok {
			user.Email = email
			return user, nil
		}
	}

	return user, errors.New("invalid token")
}

// Валидация refresh токена
func (j *JWT) ValidateRefreshToken(refreshToken string) (*user.User, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.RefreshTokenSecretKey), nil
	})

	user := &user.User{}
	if err != nil {
		return user, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if tokenType, ok := claims["type"].(string); !ok || tokenType != "refresh" {
			return user, errors.New("invalid token type")
		}

		if email, ok := claims["user_id"].(string); ok {
			user.Email = email
			return user, nil
		}
	}

	return user, errors.New("invalid token")
}

// Обновление токенов по refresh токену
func (j *JWT) RefreshTokens(refreshToken string) (*user.User, string, string, error) {
	user, err := j.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, "", "", err
	}

	accessToken, refreshToken, err := j.CreateTokenPair(user)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}
