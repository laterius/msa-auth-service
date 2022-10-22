package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/laterius/service_architecture_hw3/app/internal/domain"
	"github.com/laterius/service_architecture_hw3/app/internal/service"
)

func NewDeleteUser(d service.UserDeleter) *deleteUserHandler {
	return &deleteUserHandler{
		deleter: d,
	}
}

type deleteUserHandler struct {
	deleter service.UserDeleter
}

func (h *deleteUserHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, err := ctx.ParamsInt(UserIdFieldName, 0)
		if err != nil {
			return fail(ctx, err)
		}

		err = h.deleter.Delete(domain.UserId(userId))
		if err != nil {
			return fail(ctx, err)
		}
		return ctx.SendStatus(fiber.StatusNoContent)
	}
}
