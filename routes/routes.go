// routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/toggle-assignment/handlers"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()
	router.GET("/receipt/:uuid", handlers.GetReceipt)
	router.POST("/receipt", handlers.UploadReceipt)
	return router
}
