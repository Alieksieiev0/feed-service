package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Alieksieiev0/feed-service/api/proto"
	"github.com/Alieksieiev0/feed-service/internal/models"
	"github.com/Alieksieiev0/feed-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "No authorization code"})
		}

		claims, err := cfg.Client.ReadClaims(c.Context(), token)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		user := &models.User{
			Base: models.Base{
				ID: claims.UserId,
			},
			Name: claims.Username,
		}
		c.Locals("user", user)

		return c.Next()
	}
}

type UserConfig struct {
	Serv services.UserService
}

func NewUserMiddleware(cfg UserConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := getUserFromLocals(c)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		feedUser, err := cfg.Serv.Get(c.Context(), user.ID)
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		err = cfg.Serv.Save(c.Context(), feedUser)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Next()
	}
}

func getUserFromLocals(c *fiber.Ctx) (*models.User, error) {
	user := c.Locals("user")
	if user == nil {
		return nil, fmt.Errorf("user was not found in locals")
	}

	convUser, ok := user.(*models.User)
	if !ok {
		return nil, fmt.Errorf("invalid user data in locals")
	}
	return convUser, nil
}
