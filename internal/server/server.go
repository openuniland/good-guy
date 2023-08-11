package server

import (
	"github.com/gin-gonic/gin"
	"github.com/openuniland/good-guy/configs"
	"github.com/openuniland/good-guy/pkg/db/mongodb"
	utils "github.com/openuniland/good-guy/pkg/logger"
	"github.com/rs/zerolog/log"
)

type Server struct {
	configs *configs.Configs
	router  *gin.Engine
	mongo   *mongodb.MongoDB
}

func NewServer(configs *configs.Configs, mongo *mongodb.MongoDB) (*Server, error) {

	server := &Server{configs: configs, mongo: mongo}

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
