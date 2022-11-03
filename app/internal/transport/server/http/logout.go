package http

import (
	"github.com/gofiber/fiber/v2"
)

func Logout() *logoutGetHandler {
	return &logoutGetHandler{}
}

type logoutGetHandler struct {
}

func (h *logoutGetHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.ClearCookie("remember_token")
		return ctx.SendString("Logout success")
	}
}
