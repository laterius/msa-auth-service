package http

import (
	"github.com/gofiber/fiber/v2"
)

func GetHomePage() *getHomePageHandler {
	return &getHomePageHandler{}
}

type getHomePageHandler struct {
}

func (h *getHomePageHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Render("home", fiber.Map{})
	}
}
