package httpmw

import (
	"github.com/gofiber/fiber/v2"
	"math"
	"math/rand"
	"time"
)

const FailureChance = 10
const MaxDelayMs = 500

func NewChaosMonkeyMw() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if rand.Intn(99) < FailureChance {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}

		d := MaxDelayMs * math.Sin(float64(rand.Intn(MaxDelayMs))/float64(MaxDelayMs)*math.Pi)
		time.Sleep(time.Duration(d) * time.Millisecond)

		return ctx.Next()
	}
}
