package server

import (
	"github.com/gin-gonic/gin"
	http "github.com/openuniland/good-guy/external/ctms/delivery"
	"github.com/openuniland/good-guy/external/ctms/usecase"
	cors "github.com/rs/cors/wrapper/gin"
)

func (server *Server) MapHandlers() {
	router := gin.Default()
	router.Use(cors.AllowAll())

	// Init repositories

	// Init useCases
	ctmsUC := usecase.NewCtmsUseCase(server.configs)

	// Init handlers
	authHandlers := http.NewCtmsHandlers(server.configs, ctmsUC)

	// Init routes
	v1 := router.Group("/v1")

	health := v1.Group("/health")
	ctms := v1.Group("/ctms")

	health.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Map routes
	http.MapCtmsRoutes(ctms, authHandlers)

	server.router = router
}
