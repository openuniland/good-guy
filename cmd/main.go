package main

import (
	"fmt"
	"os"
	"time"

	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/internal/server"
	"github.com/openuniland/good-guy/pkg/db/mongodb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "time/tzdata"
)

func main() {
	configs, err := configs.LoadConfigs(".")
	if err != nil {
		fmt.Println(err)
	}

	if configs.Server.Env == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msgf("load configs success")

	// Set the default time zone to "Asia/Ho_Chi_Minh"
	loc := time.FixedZone("GMT+7", 7*60*60)
	time.Local = loc
	currentTime := time.Now()
	fmt.Println("Current time in Vietnam:", currentTime.Format(time.RFC3339))

	mongo, err := mongodb.NewMongoDB(configs)

	if err != nil {
		log.Fatal().Err(err).Msg("cannot create mongodb client")
	}
	defer mongo.Close()

	runGinServer(configs, mongo)

}

func runGinServer(config *configs.Configs, mongo *mongodb.MongoDB) {
	server, err := server.NewServer(config, mongo)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}

	err = server.Start(config.Server.HttpServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
