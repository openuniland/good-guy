package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ctmsHttp "github.com/openuniland/good-guy/external/ctms/delivery"
	facebookHttp "github.com/openuniland/good-guy/external/facebook/delivery"
	fithouHttp "github.com/openuniland/good-guy/external/fithou/delivery"
	articleHttp "github.com/openuniland/good-guy/internal/articles/delivery"
	authHttp "github.com/openuniland/good-guy/internal/auth/delivery"
	commonHttp "github.com/openuniland/good-guy/internal/common/delivery"
	cookieHttp "github.com/openuniland/good-guy/internal/cookies/delivery"
	examSchedulesHttp "github.com/openuniland/good-guy/internal/exam_schedules/delivery"
	"github.com/openuniland/good-guy/internal/middlewares"
	userHttp "github.com/openuniland/good-guy/internal/users/delivery"
	"github.com/openuniland/good-guy/jobs"

	ctmsUS "github.com/openuniland/good-guy/external/ctms/usecase"
	facebookUS "github.com/openuniland/good-guy/external/facebook/usecase"
	fithouUS "github.com/openuniland/good-guy/external/fithou/usecase"
	articleUS "github.com/openuniland/good-guy/internal/articles/usecase"
	authUC "github.com/openuniland/good-guy/internal/auth/usecase"
	commonUC "github.com/openuniland/good-guy/internal/common/usecase"
	cookieUS "github.com/openuniland/good-guy/internal/cookies/usecase"
	examSchedulesUS "github.com/openuniland/good-guy/internal/exam_schedules/usecase"
	userUS "github.com/openuniland/good-guy/internal/users/usecase"

	articleRepo "github.com/openuniland/good-guy/internal/articles/repository"
	cookieRepo "github.com/openuniland/good-guy/internal/cookies/repository"
	examSchedulesRepo "github.com/openuniland/good-guy/internal/exam_schedules/repository"
	userRepo "github.com/openuniland/good-guy/internal/users/repository"
	cors "github.com/rs/cors/wrapper/gin"

	_ "embed"
)

func (server *Server) MapHandlers() {
	router := gin.Default()
	router.Use(cors.AllowAll())

	router.Static("/static", "./static")
	router.LoadHTMLGlob("templates/*")

	// Init repositories
	articleRepo := articleRepo.NewArticleRepository(server.configs, server.mongo)
	userRepo := userRepo.NewUserRepository(server.configs, server.mongo)
	examSchedulesRepo := examSchedulesRepo.NewExamSchedulesRepository(server.configs, server.mongo)
	cookieRepo := cookieRepo.NewCookieRepository(server.configs, server.mongo)

	// Init useCases
	cookieUC := cookieUS.NewCookieUseCase(server.configs, cookieRepo)
	fithouUS := fithouUS.NewFithouUseCase(server.configs)
	articleUS := articleUS.NewArticleUseCase(server.configs, articleRepo, fithouUS)
	facebookUS := facebookUS.NewFacebookUseCase(server.configs)
	examSchedulesUS := examSchedulesUS.NewExamSchedulesUseCase(server.configs, examSchedulesRepo)
	userUS := userUS.NewUserUseCase(server.configs, userRepo)
	ctmsUC := ctmsUS.NewCtmsUseCase(server.configs, examSchedulesUS, facebookUS, userUS, cookieUC)
	authUC := authUC.NewAuthUseCase(server.configs, ctmsUC, userUS, facebookUS, cookieUC)
	commonUC := commonUC.NewCommonUseCase(server.configs, facebookUS, ctmsUC, userUS, articleUS, cookieUC)

	// Init handlers
	authCtmsHandlers := ctmsHttp.NewCtmsHandlers(server.configs, ctmsUC)
	articleHandlers := articleHttp.NewArticleHandlers(server.configs, articleUS)
	fithouHandlers := fithouHttp.NewFithouHandlers(server.configs, fithouUS)
	facebookHandlers := facebookHttp.NewFacebookHandlers(server.configs, facebookUS)
	userHandlers := userHttp.NewArticleHandlers(server.configs, userUS)
	examSchedulesHandlers := examSchedulesHttp.NewExamSchedulesHandlers(server.configs, examSchedulesUS)
	authHandlers := authHttp.NewCtmsHandlers(server.configs, authUC)
	commonHandlers := commonHttp.NewCommonHandlers(server.configs, commonUC)
	cookieHandlers := cookieHttp.NewCookieHandlers(server.configs, cookieUC)

	// Init middlewares
	mw := middlewares.NewMiddlewareManager(server.configs)

	// Jobs
	jobs := jobs.NewJobs(server.configs, articleUS, userUS, facebookUS, ctmsUC)
	jobs.Run()

	// Init web

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/privacy-policy", func(c *gin.Context) {
		c.HTML(http.StatusOK, "privacy-policy.html", nil)
	})

	// Init routes
	v1 := router.Group("/v1")

	health := v1.Group("/health")
	ctms := v1.Group("/ctms")
	articles := v1.Group("/articles")
	fithou := v1.Group("/fithou")
	facebook := v1.Group("/facebook")
	users := v1.Group("/users")
	examSchedules := v1.Group("/exam-schedules")
	auth := v1.Group("/auth")
	common := v1.Group("/common")
	cookie := v1.Group("/cookies")

	health.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Map routes
	ctmsHttp.MapCtmsRoutes(ctms, authCtmsHandlers, mw)
	articleHttp.MapArticleRoutes(articles, articleHandlers, mw)
	fithouHttp.MapFithouRoutes(fithou, fithouHandlers, mw)
	facebookHttp.MapFacebookRoutes(facebook, facebookHandlers, mw)
	userHttp.MapUserRoutes(users, userHandlers, mw)
	examSchedulesHttp.MapExamSchedulesRoutes(examSchedules, examSchedulesHandlers, mw)
	commonHttp.MapCommonRoutes(common, commonHandlers)
	authHttp.MapAuthRoutes(auth, authHandlers)
	cookieHttp.MapCookieRoutes(cookie, cookieHandlers, mw)

	server.router = router
}
