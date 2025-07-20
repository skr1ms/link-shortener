package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"not null;uniqueIndex:idx_email"`
	Password string `gorm:"not null"`
	Name     string `gorm:"not null"`
}

func NewUser(email, password, name string) *User {
	return &User{
		Email:    email,
		Password: password,
		Name:     name,
	}
}