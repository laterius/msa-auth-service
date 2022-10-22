package http

import (
	"github.com/gofiber/fiber/v2"
)

func RespondOk(ctx *fiber.Ctx) error {
	return json(ctx, struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	})
}
