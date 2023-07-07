package server

import (
	"github.com/gin-gonic/gin"
	ctmsHttp "github.com/openuniland/good-guy/external/ctms/delivery"
	fithouHttp "github.com/openuniland/good-guy/external/fithou/delivery"
	articleHttp "github.com/openuniland/good-guy/internal/articles/delivery"

	ctmsUS "github.com/openuniland/good-guy/external/ctms/usecase"
	fithouUS "github.com/openuniland/good-guy/external/fithou/usecase"
	articleUS "github.com/openuniland/good-guy/internal/articles/usecase"

	articleRepo "github.com/openuniland/good-guy/internal/articles/repository"
	cors "github.com/rs/cors/wrapper/gin"
)

func (server *Server) MapHandlers() {
	router := gin.Default()
	router.Use(cors.AllowAll())

	// Init repositories
	articleRepo := articleRepo.NewArticleRepository(server.configs, server.mongoClient)

	// Init useCases
	ctmsUC := ctmsUS.NewCtmsUseCase(server.configs)
	fithouUS := fithouUS.NewFithouUseCase(server.configs)
	articleUS := articleUS.NewArticleUseCase(server.configs, articleRepo, fithouUS)

	// Init handlers
	authHandlers := ctmsHttp.NewCtmsHandlers(server.configs, ctmsUC)
	articleHandlers := articleHttp.NewArticleHandlers(server.configs, articleUS)
	fithouHandlers := fithouHttp.NewFithouHandlers(server.configs, fithouUS)

	// Init routes
	v1 := router.Group("/v1")

	health := v1.Group("/health")
	ctms := v1.Group("/ctms")
	articles := v1.Group("/articles")
	fithou := v1.Group("/fithou")

	health.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Map routes
	ctmsHttp.MapCtmsRoutes(ctms, authHandlers)
	articleHttp.MapArticleRoutes(articles, articleHandlers)
	fithouHttp.MapFithouRoutes(fithou, fithouHandlers)

	server.router = router
}
