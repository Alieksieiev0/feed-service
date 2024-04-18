package rest

import (
	"fmt"
	"strconv"

	"github.com/Alieksieiev0/feed-service/internal/models"
	"github.com/Alieksieiev0/feed-service/internal/services"
	"github.com/Alieksieiev0/feed-service/internal/transport/kafka"
	"github.com/Alieksieiev0/feed-service/internal/types"
	"github.com/gofiber/fiber/v2"
)

func GetPosts(serv services.FeedService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params, err := getDefaultParams(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		posts, err := serv.GetPosts(c.Context(), params...)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(posts)
	}
}

func GetUsers(serv services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		params, err := getDefaultParams(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		username := c.Query("username")
		params = append(params, services.Filter("username", username, false))

		users, err := serv.GetUsers(c.Context(), params...)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		for _, u := range users {
			u.Password = ""
		}

		return c.Status(fiber.StatusOK).JSON(users)
	}
}

func Subscribe(serv services.UserFeedService, producer kafka.Producer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sub := &models.User{}
		if err := c.BodyParser(sub); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		userId := c.Params("id")
		user, err := serv.GetById(c.Context(), userId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err = serv.Subscribe(c.Context(), user.Id, sub.ID); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		err = producer.Produce(
			[]types.UserBase{user.UserBase},
			kafka.NewSubscriptionNotification(sub.ID, sub.Username),
		)
		fmt.Println("-----")
		fmt.Println(err)
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

		c.Status(fiber.StatusOK)
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
		user, err := serv.GetById(c.Context(), userId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err = serv.Post(c.Context(), user.Id, post); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		err = producer.Produce(
			user.Subscribers,
			kafka.NewPostNotification(user.Id, user.Username, post.ID),
		)
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

		return c.Status(fiber.StatusCreated).JSON(post)
	}
}

func getDefaultParams(c *fiber.Ctx) ([]services.Param, error) {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil {
		return nil, fmt.Errorf("invalid limit was provided")
	}

	offset, err := strconv.Atoi(c.Query("offset", "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid offset was provided")
	}

	sortBy := c.Query("sort_by", "Id")
	orderBy := c.Query("order_by", "asc")

	params := []services.Param{
		services.Limit(limit),
		services.Offset(offset),
		services.Order(sortBy, orderBy),
	}
	return params, nil
}
