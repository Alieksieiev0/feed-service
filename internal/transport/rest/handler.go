package rest

import (
	"fmt"
	"log"
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

		if n := c.Query("username"); n != "" {
			params = append(params, services.Filter("username", n, false))
		}

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

func Subscribe(serv services.FeedService, producer kafka.Producer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sub := &models.User{}
		if err := c.BodyParser(sub); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		userId := c.Params("id")
		if err := serv.Subscribe(c.Context(), userId, sub.ID); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		go func() {
			err := producer.Produce(
				[]types.UserBase{{Id: userId}},
				kafka.NewSubscriptionMessage(sub.ID, sub.Username),
			)
			if err != nil {
				log.Println(err)
			}
		}()

		c.Status(fiber.StatusOK)
		return nil
	}
}

func Unsubscribe(serv services.FeedService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sub := &models.User{}
		if err := c.BodyParser(sub); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := serv.Unsubscribe(c.Context(), c.Params("id"), sub.ID); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
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

		user, err := serv.GetById(c.Context(), c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		if err = serv.Post(c.Context(), user.Id, post); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
		}

		go func() {
			err := producer.Produce(
				user.Subscribers,
				kafka.NewPostMessage(user.Id, user.Username, post.ID),
			)
			if err != nil {
				log.Println(err)
			}
		}()

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
