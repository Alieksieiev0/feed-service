package rest

import (
	"github.com/Alieksieiev0/feed-service/api/proto"
	"github.com/Alieksieiev0/feed-service/internal/services"
	"github.com/Alieksieiev0/feed-service/internal/transport/kafka"
	"github.com/gofiber/fiber/v2"
)

type RESTServer struct {
	app *fiber.App
}

func NewServer(app *fiber.App) *RESTServer {
	return &RESTServer{
		app: app,
	}
}

func (us *RESTServer) Start(
	addr string,
	serv services.UserFeedService,
	client proto.AuthServiceClient,
	producer kafka.Producer,
) error {
	us.app.Use(
		NewAuthMiddleware(AuthConfig{Client: client}),
		NewUserMiddleware(UserConfig{Serv: serv}),
	)

	us.app.Put("/users/:id/subscribers", Subscribe(serv, producer))
	us.app.Put("/users/:id/posts", Post(serv, producer))

	return us.app.Listen(addr)
}