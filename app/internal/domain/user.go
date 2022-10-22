package domain

import "github.com/laterius/service_architecture_hw3/app/pkg/types"

type UserId int64
type Username string
type Password string

func (t UserId) Validate() error {
	if t <= 0 {
		return ErrInvalidUserId
	}
	return nil
}

type User struct {
	Id           UserId
	Username     string
	FirstName    string
	LastName     string
	Email        string
	Phone        string
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

type UserPartialData = types.Kv
