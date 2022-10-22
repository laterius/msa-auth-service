package http

import "github.com/gofiber/fiber/v2"

func SignUpGet() *signUpGetHandler {
	return &signUpGetHandler{}
}

type signUpGetHandler struct {
}

func (h *signUpGetHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Render("signup", fiber.Map{})
	}
}
