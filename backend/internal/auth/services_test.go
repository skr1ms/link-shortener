package auth_test

import (
	"errors"
	"linkshortener/internal/auth"
	"linkshortener/internal/user"
	"linkshortener/pkg/jwt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	users map[string]*user.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*user.User),
	}
}

func (m *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *MockUserRepository) Create(u *user.User) (*user.User, error) {
	if _, exists := m.users[u.Email]; exists {
		return nil, errors.New("user already exists")
	}
	m.users[u.Email] = u
	return u, nil
}

func setupAuthService() (*auth.AuthService, *MockUserRepository) {
	mockRepo := NewMockUserRepository()
	jwtService := jwt.NewJWT(os.Getenv("SECRET_KEY"), os.Getenv("REFRESH_SECRET_KEY"))
	authService := auth.NewAuthService(mockRepo, jwtService)
	return authService, mockRepo
}

func TestAuthServiceRegisterSuccess(t *testing.T) {
	godotenv.Load()
	authService, _ := setupAuthService()

	user, err := authService.Register("test@example.com", "password123", "Test User")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("Expected user to be created")
	}

	if user.Email != "test@example.com" {
		t.Fatalf("Expected email test@example.com, got %s", user.Email)
	}

	if user.Name != "Test User" {
		t.Fatalf("Expected name Test User, got %s", user.Name)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	if err != nil {
		t.Fatal("Password was not hashed correctly")
	}
}

func TestAuthServiceRegisterUserAlreadyExists(t *testing.T) {
	godotenv.Load()
	authService, mockRepo := setupAuthService()

	existingUser := &user.User{
		Email: "test@example.com",
		Name:  "Existing User",
	}
	mockRepo.users["test@example.com"] = existingUser

	_, err := authService.Register("test@example.com", "password123", "Test User")

	if err == nil {
		t.Fatal("Expected error for existing user")
	}

	if err.Error() != "user already exists" {
		t.Fatalf("Expected 'user already exists' error, got %v", err)
	}
}

func TestAuthServiceLoginSuccess(t *testing.T) {
	godotenv.Load()
	authService, mockRepo := setupAuthService()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	existingUser := &user.User{
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Name:     "Test User",
	}
	mockRepo.users["test@example.com"] = existingUser

	user, err := authService.Login("test@example.com", "password123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user == nil {
		t.Fatal("Expected user to be returned")
	}

	if user.Email != "test@example.com" {
		t.Fatalf("Expected email test@example.com, got %s", user.Email)
	}
}

func TestAuthServiceLoginUserNotFound(t *testing.T) {
	godotenv.Load()
	authService, _ := setupAuthService()

	_, err := authService.Login("nonexistent@example.com", "password123")

	if err == nil {
		t.Fatal("Expected error for non-existent user")
	}

	if err.Error() != "user not found" {
		t.Fatalf("Expected 'user not found' error, got %v", err)
	}
}

func TestAuthServiceLoginInvalidPassword(t *testing.T) {
	godotenv.Load()
	authService, mockRepo := setupAuthService()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)
	existingUser := &user.User{
		Email:    "test@example.com",
		Password: string(hashedPassword),
		Name:     "Test User",
	}
	mockRepo.users["test@example.com"] = existingUser

	_, err := authService.Login("test@example.com", "wrong_password")

	if err == nil {
		t.Fatal("Expected error for wrong password")
	}

	if err.Error() != "invalid password" {
		t.Fatalf("Expected 'invalid password' error, got %v", err)
	}
}
