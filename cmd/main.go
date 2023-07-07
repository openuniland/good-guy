package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/internal/server"
	"github.com/openuniland/good-guy/pkg/db/mongodb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {

	configs, err := configs.LoadConfigs(".")
	if err != nil {
		fmt.Println(err)
	}

	if configs.Server.Env == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	redisClient, err := mongodb.NewMongoDBClient(configs)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create mongodb client")
	}
	defer redisClient.Disconnect(context.Background())

	runGinServer(configs, redisClient)

}

func runGinServer(config *configs.Configs, mongoClient *mongo.Client) {
	server, err := server.NewServer(config, mongoClient)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.Server.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
