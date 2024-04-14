package rest

import (
	"github.com/Alieksieiev0/feed-service/api/proto"
	"github.com/gofiber/fiber/v2"
)

type AuthConfig struct {
	Client proto.AuthServiceClient
}

func NewAuthMiddleware(cfg AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := &proto.Token{
			Value: c.Get("Authorization"),
		}
		if token.Value == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(fiber.Map{"error": "No authorization code"})
		}

		_, err := cfg.Client.ReadClaims(c.Context(), token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Next()
	}
}
