package di

import (
	"linkshortener/internal/user"
	"linkshortener/pkg/event"
)

type IEventBus interface {
	Publish(event event.Event)
	Subscribe() <-chan event.Event
}

type IStatsRepository interface {
	AddClick(linkId uint) error
}

type IUserRepository interface {
	Create(user *user.User) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
}
