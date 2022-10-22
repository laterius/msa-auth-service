package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/laterius/service_architecture_hw3/app/internal/service"
)

func NewGetProfile(r service.UserRememberReader) *getProfileHandler {
	return &getProfileHandler{
		reader: r,
	}
}

type getProfileHandler struct {
	reader service.UserRememberReader
}

func (h *getProfileHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		//TODO доделать
		token := ctx.Params(TokenFieldName, "")
		if token == "" {
			return fail(ctx, errors.New("token is empty"))
		}

		user, err := h.reader.ByRemember(token)
		if err != nil {
			return fail(ctx, err)
		}

		return ctx.Render("profile", fiber.Map{
			"FirstName": user.FirstName,
			"LastName":  user.LastName,
			"Username":  user.Username,
			"Phone":     user.Phone,
			"Email":     user.Email,
			"Token":     user.Remember,
		})
	}
}
