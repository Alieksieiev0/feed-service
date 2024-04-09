package rest

import (
	"net/http"

	"github.com/Alieksieiev0/feed-service/internal/models"
	"github.com/Alieksieiev0/feed-service/internal/services"
	"github.com/Alieksieiev0/feed-service/internal/transport/kafka"
	"github.com/Alieksieiev0/feed-service/internal/types"
	"github.com/gofiber/fiber/v2"
)

func Subscribe(serv services.UserFeedService, producer kafka.Producer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sub := &models.User{}
		if err := c.BodyParser(sub); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		userId := c.Params("id")
		user, err := serv.Get(c.Context(), userId)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err = serv.Subscribe(c.Context(), user, sub); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		err = producer.Produce([]models.User{*user}, kafka.NewSubscription(sub.ID, sub.Name))
		if err != nil {
			body := types.SubscriptionPartialSuccess{
				Subscription: types.XMLResponse{
					Status: fiber.StatusOK,
				},
				Notification: types.XMLResponse{
					Status: fiber.StatusInternalServerError,
					Error:  err.Error(),
				},
			}

			return c.Status(fiber.StatusMultiStatus).XML(body)
		}

		c.Status(http.StatusOK)
		return nil
	}
}

func Post(serv services.UserFeedService, producer kafka.Producer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		post := &models.Post{}
		if err := c.BodyParser(post); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("error", err.Error())
		}

		userId := c.Params("id")
		user, err := serv.Get(c.Context(), userId)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err = serv.Post(c.Context(), user, post); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		err = producer.Produce(user.Subcribers, kafka.NewPost(user.ID, user.Name, post.ID))
		if err != nil {
			body := types.PostPartialSuccess{
				Creation: types.XMLResponse{
					Status: fiber.StatusCreated,
				},
				Notification: types.XMLResponse{
					Status: fiber.StatusInternalServerError,
					Error:  err.Error(),
				},
			}

			return c.Status(fiber.StatusMultiStatus).XML(body)
		}

		c.Status(http.StatusCreated)
		return nil
	}
}
