package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"mihiru-go/config"
	"mihiru-go/controllers"
	"mihiru-go/database"
	"mihiru-go/middleware"
	"mihiru-go/services"
)

var adminRole = []string{"admin"}

func NewRouter(db *database.MongoDatabase) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.GetConfigs().GetStringSlice("server.allow-origins")
	router.Use(cors.New(corsConfig))

	userService := services.NewUserService(db)
	userService.InitUser()
	permissions := middleware.NewPermissions(userService)

	articlesGroup := router.Group("articles")
	{
		articleService := services.NewArticleService(db)
		articlesController := controllers.NewArticlesController(articleService, userService)
		articlesGroup.POST("/", permissions.Roles(adminRole), articlesController.Add)
		articlesGroup.POST("/search", articlesController.Search)
		articlesGroup.GET("/tags", articlesController.Tags)
		articlesGroup.PUT("/:id", permissions.Roles(adminRole), articlesController.Update)
		articlesGroup.GET("/:id", articlesController.Get)
	}

	userGroup := router.Group("user")
	{
		userController := controllers.NewUserController(userService)
		userGroup.POST("/", permissions.Roles(adminRole), userController.Add)
		userGroup.POST("/login", userController.Login)
	}

	return router
}
