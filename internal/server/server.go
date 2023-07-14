package server

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	utils "github.com/openuniland/good-guy/pkg/logger"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	configs     *configs.Configs
	router      *gin.Engine
	mongoClient *mongo.Client
}

func NewServer(configs *configs.Configs, mongoClient *mongo.Client) (*Server, error) {

	server := &Server{configs: configs, mongoClient: mongoClient}

	if server.configs.Server.Env != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	return server, nil
}

func (server *Server) Start(address string) error {

	server.MapHandlers()
	log.Info().Msg("starting HTTP server")
	if server.configs.Server.Env != "dev" {
		return server.router.Run()
	}
	return server.router.Run(":" + address)
}

func (server *Server) HttpLogger() gin.IRoutes {
	return server.router.Use(utils.HttpLogger())
}
