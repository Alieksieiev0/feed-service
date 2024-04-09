package main

import (
	"flag"
	"log"

	"github.com/Alieksieiev0/feed-service/internal/database"
	"github.com/Alieksieiev0/feed-service/internal/services"
	"github.com/Alieksieiev0/feed-service/internal/transport/grpc"
	"github.com/Alieksieiev0/feed-service/internal/transport/kafka"
	"github.com/Alieksieiev0/feed-service/internal/transport/rest"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	var (
		restServerAddr = flag.String("rest-server", ":3002", "listen address of rest server")
		grpcClientAddr = flag.String("grpc-client", ":4001", "listen address of grpc client")
		kafkaAddr      = flag.String("kafka", "9092", "address of kafka")
		producer       = kafka.NewProducer(*kafkaAddr)
		app            = fiber.New()
	)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Start()
	if err != nil {
		log.Fatal(err)
	}

	client, err := grpc.NewGRPCClient(*grpcClientAddr)
	if err != nil {
		log.Fatal(err)
	}

	userFeedService := services.NewUserFeedService(db)

	restServer := rest.NewServer(app)
	err = restServer.Start(*restServerAddr, userFeedService, client, producer)
	if err != nil {
		log.Fatal(err)
	}
}
