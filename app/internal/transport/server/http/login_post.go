package http

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/laterius/service_architecture_hw3/app/internal/domain"
	"github.com/laterius/service_architecture_hw3/app/internal/service"
)

func LoginPost(c service.UserLoginReader) *loginPostHandler {
	return &loginPostHandler{
		login: c,
	}
}

type loginPostHandler struct {
	login service.UserLoginReader
}

func (h *loginPostHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var u service.UserLogin
		err := ctx.BodyParser(&u)
		if err != nil {
			return fail(ctx, err)
		}

		user, err := h.login.Login(domain.Username(u.Username), domain.Password(u.Password))
		if err != nil {
			return fail(ctx, err)
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     "remember_token",
			Value:    user.Remember,
			HTTPOnly: true,
		})

		return ctx.SendString(fmt.Sprintf("Login success! login = %s, token = %s", user.Username, user.Remember))
	}
}
