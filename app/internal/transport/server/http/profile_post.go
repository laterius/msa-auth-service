package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/laterius/service_architecture_hw3/app/internal/domain"
	"github.com/laterius/service_architecture_hw3/app/internal/service"
)

func NewPostProfile(r service.UserRememberReader, pu service.UserPartialUpdater) *postProfileHandler {
	return &postProfileHandler{
		reader:         r,
		partialUpdater: pu,
	}
}

type postProfileHandler struct {
	reader         service.UserRememberReader
	partialUpdater service.UserPartialUpdater
}

func (h *postProfileHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		userId, err := ctx.ParamsInt(UserIdFieldName)
		if err != nil {
			return fail(ctx, err)
		}

		var u service.UserPartialUpdate
		err = ctx.BodyParser(&u)
		if err != nil {
			return fail(ctx, err)
		}

		currentUser, err := h.reader.ByRemember(u.Remember.Value)
		if err != nil {
			return fail(ctx, err)
		}

		if int(currentUser.Id) != userId {
			return fail(ctx, errors.New("id is empty"))
		}

		updatedUser, err := h.partialUpdater.PartialUpdate(domain.UserId(userId), u.ToDomain())
		if err != nil {
			return fail(ctx, err)
		}
		return ctx.Render("profile", fiber.Map{
			"FirstName": updatedUser.FirstName,
			"LastName":  updatedUser.LastName,
			"Username":  updatedUser.Username,
			"Phone":     updatedUser.Phone,
			"Email":     updatedUser.Email,
			"Token":     updatedUser.Remember,
		})
	}
}
