package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toggle-assignment/data"
)

type errorResponse struct {
	Code    int
	Message string
}

func GetReceipt(c *gin.Context) {
	uuid := c.Param("uuid")

	if uuid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	database := data.GetDatabase()

	receipt := database.GetReceipt(uuid)

	if receipt == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "receipt not found"})
		return
	}

	if receipt.IsProcessing {
		c.JSON(http.StatusOK, gin.H{"message": "receipt is still processing"})
		return
	}

	c.IndentedJSON(http.StatusOK, receipt)
}

func UploadReceipt(c *gin.Context) {
	defer handlePanic(c)

	file, err := c.FormFile("file")

	validateReceipt(file, err)

	newUUID := uuid.New()

	var uuid = newUUID.String()

	filePath := saveFile(c, file, uuid)

	database := data.GetDatabase()

	database.AddReceipt(uuid)

	go processReceipt(uuid, filePath)

	c.IndentedJSON(http.StatusOK, uuid)
}
