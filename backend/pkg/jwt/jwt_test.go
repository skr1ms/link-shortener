package jwt_test

import (
	"fmt"
	"os"
	"testing"

	"linkshortener/internal/user"
	"linkshortener/pkg/jwt"

	"github.com/joho/godotenv"
)

func TestGenerateToken(t *testing.T) {
	godotenv.Load()
	token, refreshToken, err := jwt.NewJWT(os.Getenv("SECRET_KEY"), os.Getenv("REFRESH_SECRET_KEY")).
		CreateTokenPair(&user.User{
			Email: "test@test.com",
		})
	if err != nil {
		t.Fatalf("Error generating token: %v", err)
	}

	if token == "" || refreshToken == "" {
		t.Fatalf("expected token and refresh token, got %s and %s", token, refreshToken)
	}

	user, err := jwt.NewJWT(os.Getenv("SECRET_KEY"), os.Getenv("REFRESH_SECRET_KEY")).ValidateToken(token)
	if err != nil {
		t.Fatalf("Error validating token: %v", err)
	}

	if user.Email != "test@test.com" {
		t.Fatalf("expected user email, got %s", user.Email)
	}

	user, err = jwt.NewJWT(os.Getenv("SECRET_KEY"), os.Getenv("REFRESH_SECRET_KEY")).ValidateRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("Error validating refresh token: %v", err)
	}

	if user.Email != "test@test.com" {
		t.Fatalf("expected user email, got %s", user.Email)
	}
	fmt.Println(user.Email)
	fmt.Println(token)
	fmt.Println(refreshToken)
}
