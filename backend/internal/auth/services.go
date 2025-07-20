package auth

import (
	"errors"
	"linkshortener/internal/user"
	"linkshortener/pkg/di"
	"linkshortener/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository di.IUserRepository
	jwt            *jwt.JWT
}

func NewAuthService(userRepository di.IUserRepository, jwt *jwt.JWT) *AuthService {
	return &AuthService{userRepository: userRepository, jwt: jwt}
}

func (service *AuthService) Register(email, password, name string) (*user.User, error) {
	existingUser, err := service.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := user.NewUser(email, string(hashedPassword), name)

	_, err = service.userRepository.Create(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (service *AuthService) Login(email, password string) (*user.User, error) {
	userExists, err := service.userRepository.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if userExists == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userExists.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return userExists, nil
}
