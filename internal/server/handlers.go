package server

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	ctmsHttp "github.com/openuniland/good-guy/external/ctms/delivery"
	facebookHttp "github.com/openuniland/good-guy/external/facebook/delivery"
	fithouHttp "github.com/openuniland/good-guy/external/fithou/delivery"
	articleHttp "github.com/openuniland/good-guy/internal/articles/delivery"
	authHttp "github.com/openuniland/good-guy/internal/auth/delivery"
	commonHttp "github.com/openuniland/good-guy/internal/common/delivery"
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
	examSchedulesUS "github.com/openuniland/good-guy/internal/exam_schedules/usecase"
	userUS "github.com/openuniland/good-guy/internal/users/usecase"

	articleRepo "github.com/openuniland/good-guy/internal/articles/repository"
	examSchedulesRepo "github.com/openuniland/good-guy/internal/exam_schedules/repository"
	userRepo "github.com/openuniland/good-guy/internal/users/repository"
	cors "github.com/rs/cors/wrapper/gin"

	_ "embed"
)

func (server *Server) MapHandlers() {
	router := gin.Default()
	router.Use(cors.AllowAll())

	router.Static("/static", "./static")

	// Init repositories
	articleRepo := articleRepo.NewArticleRepository(server.configs, server.mongo)
	userRepo := userRepo.NewUserRepository(server.configs, server.mongo)
	examSchedulesRepo := examSchedulesRepo.NewExamSchedulesRepository(server.configs, server.mongo)

	// Init useCases
	fithouUS := fithouUS.NewFithouUseCase(server.configs)
	articleUS := articleUS.NewArticleUseCase(server.configs, articleRepo, fithouUS)
	facebookUS := facebookUS.NewFacebookUseCase(server.configs)
	userUS := userUS.NewUserUseCase(server.configs, userRepo)
	examSchedulesUS := examSchedulesUS.NewExamSchedulesUseCase(server.configs, examSchedulesRepo)
	ctmsUC := ctmsUS.NewCtmsUseCase(server.configs, examSchedulesUS, facebookUS)
	authUC := authUC.NewAuthUseCase(server.configs, ctmsUC, userUS, facebookUS)
	commonUC := commonUC.NewCommonUseCase(server.configs, facebookUS, ctmsUC, userUS, articleUS)

	// Init handlers
	authCtmsHandlers := ctmsHttp.NewCtmsHandlers(server.configs, ctmsUC)
	articleHandlers := articleHttp.NewArticleHandlers(server.configs, articleUS)
	fithouHandlers := fithouHttp.NewFithouHandlers(server.configs, fithouUS)
	facebookHandlers := facebookHttp.NewFacebookHandlers(server.configs, facebookUS)
	userHandlers := userHttp.NewArticleHandlers(server.configs, userUS)
	examSchedulesHandlers := examSchedulesHttp.NewExamSchedulesHandlers(server.configs, examSchedulesUS)
	authHandlers := authHttp.NewCtmsHandlers(server.configs, authUC)
	commonHandlers := commonHttp.NewCommonHandlers(server.configs, commonUC)

	// Init middlewares
	mw := middlewares.NewMiddlewareManager(server.configs)

	// Jobs
	jobs := jobs.NewJobs(server.configs, articleUS, userUS, facebookUS, ctmsUC)
	jobs.Run()

	// Init web

	router.GET("/", func(c *gin.Context) {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Not found currentDir")
			return
		}
		fmt.Println("currentDir", currentDir)
		filePath := filepath.Join(currentDir, "/internal/server/index.html")

		file, err := os.Open(filePath)
		if err != nil {
			fmt.Printf("Không thể mở tệp: %v\n", err)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("Không thể đọc tệp: %v\n", err)
			return
		}

		fmt.Println("Nội dung tệp HTML:")
		fmt.Println(string(content))

		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Writer.WriteString(string(content))
	})

	// router.GET("/privacy-policy", func(c *gin.Context) {
	// 	t, err := template.New("privacy-policy.html").ParseFiles(frameworks.Web().PrivacyPolicy)
	// 	if err != nil {
	// 		c.JSON(500, gin.H{
	// 			"message": "Error",
	// 			"err":     err,
	// 		})
	// 		return
	// 	}

	// 	t.Execute(c.Writer, nil)
	// 	if err != nil {
	// 		c.JSON(500, gin.H{
	// 			"message": "Error",
	// 			"err":     err,
	// 		})
	// 		return
	// 	}
	// })
	router.GET("/test", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Writer.WriteString(`
		<h1>HELLO</h1>	
		`)
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

	server.router = router
}
