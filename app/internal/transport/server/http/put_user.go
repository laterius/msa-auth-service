package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/laterius/service_architecture_hw3/app/internal/domain"
	"github.com/laterius/service_architecture_hw3/app/internal/service"
)

func NewPutUser(u service.UserUpdater) *putUserHandler {
	return &putUserHandler{
		updater: u,
	}
}

type putUserHandler struct {
	updater service.UserUpdater
}

func (h *putUserHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, err := ctx.ParamsInt(UserIdFieldName, 0)
		if err != nil {
			return fail(ctx, err)
		}

		var u service.UserUpdate
		err = ctx.BodyParser(&u)
		if err != nil {
			return fail(ctx, err)
		}

		user, err := h.updater.Update(domain.UserId(userId), &u)
		if err != nil {
			return fail(ctx, err)
		}
		return json(ctx, (&service.User{}).FromDomain(user))
	}
}
