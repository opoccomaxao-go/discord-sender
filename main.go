package main

import (
	"context"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/opoccomaxao-go/discord-sender/sender"
	"github.com/opoccomaxao-go/task-server/storage"
)

type MongoConfig struct {
	ConnectURL string `env:"URL"`
	DBName     string `env:"DB_NAME"`
}

type ServerConfig struct {
	Port int `env:"PORT" envDefault:"8080"`
}

type Config struct {
	Server ServerConfig `envPrefix:"SERVER_"`
	Mongo  MongoConfig  `envPrefix:"MONGO_"`
}

func main() {
	_ = godotenv.Load(".env")

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	mongo, err := storage.NewMongo(storage.StorageMongoConfig(cfg.Mongo))
	if err != nil {
		log.Fatal(err)
	}

	sender, err := sender.New(sender.Config{
		Storage: mongo,
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := sender.Serve(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	server := Server{
		service: sender,
		config:  cfg.Server,
	}

	err = server.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
