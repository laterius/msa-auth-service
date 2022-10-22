package http

import "github.com/gofiber/fiber/v2"

func LoginGet() *loginGetHandler {
	return &loginGetHandler{}
}

type loginGetHandler struct {
}

func (h *loginGetHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Render("login", fiber.Map{})
	}
}
