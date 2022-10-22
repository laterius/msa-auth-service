package http

import (
	"github.com/gofiber/fiber/v2"
)

func GetContact() *getContactHandler {
	return &getContactHandler{}
}

type getContactHandler struct {
}

func (h *getContactHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Render("contact", fiber.Map{})
	}
}
