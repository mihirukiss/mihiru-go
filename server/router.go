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
	userController := controllers.NewUserController(userService)

	articleService := services.NewArticleService(db)
	articlesController := controllers.NewArticlesController(articleService, userService)

	memoryService := services.NewMemoryService(db, db)
	memoryController := controllers.NewMemoryController(memoryService)

	permissions := middleware.NewPermissions(userService)

	articlesGroup := router.Group("articles")
	{
		articlesGroup.POST("", permissions.Roles(adminRole), articlesController.Add)
		articlesGroup.POST("/search", articlesController.Search)
		articlesGroup.GET("/tags", articlesController.Tags)
		articlesGroup.PUT("/:id", permissions.Roles(adminRole), articlesController.Update)
		articlesGroup.GET("/:id", articlesController.Get)
	}

	userGroup := router.Group("user")
	{
		userGroup.POST("", permissions.Roles(adminRole), userController.Add)
		userGroup.POST("/login", userController.Login)
	}

	memoryGroup := router.Group("memory")
	{
		memoryGroup.POST("/dynamic", permissions.Roles(adminRole), memoryController.AddDynamic)
		memoryGroup.PUT("/dynamic/:id", permissions.Roles(adminRole), memoryController.UpdateDynamic)
		memoryGroup.POST("/live", permissions.Roles(adminRole), memoryController.AddLive)
		memoryGroup.PUT("/live/:id", permissions.Roles(adminRole), memoryController.UpdateLive)
		memoryGroup.GET("/days", memoryController.Days)
		memoryGroup.GET("/day/:day", memoryController.Day)
	}

	return router
}
