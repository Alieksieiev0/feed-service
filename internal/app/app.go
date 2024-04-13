package app

import (
	"flag"
	"log"

	"github.com/Alieksieiev0/feed-service/internal/config"
	"github.com/Alieksieiev0/feed-service/internal/services"
	"github.com/Alieksieiev0/feed-service/internal/transport/grpc"
	"github.com/Alieksieiev0/feed-service/internal/transport/kafka"
	"github.com/Alieksieiev0/feed-service/internal/transport/rest"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

func Run() {
	var (
		restServerAddr = flag.String("rest-server", ":3000", "listen address of rest server")
		grpcServerAddr = flag.String("grpc-server", ":4000", "listen address of grpc server")
		grpcClientAddr = flag.String(
			"grpc-client",
			"auth-service:4001",
			"listen address of grpc client",
		)
		kafkaAddr = flag.String("kafka", "9092", "address of kafka")
		producer  = kafka.NewProducer(*kafkaAddr)
		app       = fiber.New()
		g         = new(errgroup.Group)
	)
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := config.Database()
	if err != nil {
		log.Fatal(err)
	}

	client, err := grpc.NewGRPCClient(*grpcClientAddr)
	if err != nil {
		log.Fatal(err)
	}

	serv := services.NewUserFeedService(db)

	grpcServer := grpc.NewServer()
	g.Go(func() error {
		return grpcServer.Start(*grpcServerAddr, serv)
	})

	restServer := rest.NewServer(app, *restServerAddr)
	g.Go(func() error {
		return restServer.Start(serv, client, producer)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
