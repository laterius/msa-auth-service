package dbrepo

import (
	"fmt"
	"github.com/laterius/service_architecture_hw3/app/internal/domain"
)

func Dsn(cfg domain.Db) string {
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DbName,
		cfg.Password,
	)
	if cfg.Extras != "" {
		connString += " " + cfg.Extras
	}
	return connString
}
