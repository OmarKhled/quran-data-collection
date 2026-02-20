package routes

import (
	"data-collection/controllers"
	"data-collection/db"
	"data-collection/middleware"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db_client *db.Queries, r2_client *s3.Client) {
	var bucketName = os.Getenv("BUCKET_NAME")

	router.GET("/user", controllers.GetUser(db_client))
	router.POST("/user", controllers.CreateUser(db_client))

	router.GET("/users/rank", controllers.GetUsersRanks(db_client))
	router.GET("/user/rank", controllers.GetUserRank(db_client))

	router.GET("/task", controllers.GetTask(db_client))
	router.POST("/task", controllers.CreateUserTask(db_client, r2_client, bucketName))

	// ============= Admin routes =============
	router.POST("/admin/login", controllers.Login(db_client))
	router.POST("/admin/logout", controllers.Logout(db_client))
	adminRouter := router.Group("/admin")
	adminRouter.Use(middleware.IsAuthorized(db_client))
	adminRouter.GET("/tasks", controllers.GetUsersTasks(db_client, r2_client))
}
