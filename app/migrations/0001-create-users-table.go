package mixtures

import (
	"github.com/ezn-go/mixture"
	"github.com/go-gormigrate/gormigrate/v2"
)

func init() {
	type User struct {
		Id           int64  `json:"id" gorm:"primaryKey,autoIncrement"`
		Username     string `json:"username"`
		FirstName    string `json:"firstName"`
		LastName     string `json:"lastName"`
		Email        string `json:"email"`
		Phone        string `json:"phone"`
		Password     string `json:"password"`
		PasswordHash string `json:"passwordHash"`
		Remember     string `json:"remember"`
		RememberHash string `json:"rememberHash"`
	}

	mx := &gormigrate.Migration{
		ID:       "0001",
		Migrate:  mixture.CreateTableM(&User{}),
		Rollback: mixture.DropTableR(&User{}),
	}

	mixture.Add(mixture.ForAnyEnv, mx)
}
