package rest

import (
	"github.com/Alieksieiev0/feed-service/api/proto"
	"github.com/Alieksieiev0/feed-service/internal/services"
	"github.com/Alieksieiev0/feed-service/internal/transport/kafka"
	"github.com/gofiber/fiber/v2"
)

type RESTServer struct {
	app  *fiber.App
	addr string
}

func NewServer(app *fiber.App, addr string) *RESTServer {
	return &RESTServer{
		app:  app,
		addr: addr,
	}
}

func (us *RESTServer) Start(
	serv services.UserFeedService,
	client proto.AuthServiceClient,
	producer kafka.Producer,
) error {

	us.app.Get("/posts", GetPosts(serv))
	us.app.Use(
		NewAuthMiddleware(AuthConfig{Client: client}),
	)
	us.app.Put("/users/:id/subscribers", Subscribe(serv, producer))
	us.app.Put("/users/:id/posts", Post(serv, producer))

	return us.app.Listen(us.addr)
}
