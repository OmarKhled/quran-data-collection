package routes

import (
	"data-collection/controllers"
	"data-collection/db"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db_client *db.Queries) {
	router.GET("/user", controllers.GetUser(db_client))
	router.POST("/user", controllers.CreateUser(db_client))

	router.GET("/task", controllers.GetTask(db_client))
}
