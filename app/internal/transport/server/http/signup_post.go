package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/laterius/service_architecture_hw3/app/internal/service"
)

func SignUpPost(c service.UserCreator) *signUpPostHandler {
	return &signUpPostHandler{
		creator: c,
	}
}

type signUpPostHandler struct {
	creator service.UserCreator
}

func (h *signUpPostHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var u service.UserCreate
		err := ctx.BodyParser(&u)
		if err != nil {
			return fail(ctx, err)
		}

		user, err := h.creator.Create(&u)
		if err != nil {
			return fail(ctx, err)
		}

		ctx.Cookie(&fiber.Cookie{
			Name:     "remember_token",
			Value:    user.Remember,
			HTTPOnly: true,
		})

		return created(ctx, (&service.User{}).FromDomain(user), int64(user.Id))
	}
}
