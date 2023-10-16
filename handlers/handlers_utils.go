package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/toggle-assignment/data"
	receiptprocessor "github.com/toggle-assignment/receipt_processing"
)

func isValidFileExtension(ext string) bool {
	validExtensions := []string{".png", ".jpg", ".pdf"}
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

func handlePanic(c *gin.Context) {
	if r := recover(); r != nil {
		if response, ok := r.(errorResponse); ok {
			c.JSON(response.Code, gin.H{"error": response.Message})
		} else {
			fmt.Println("Recovered from panic with unexpected value:", r)
		}
		panic("An error occurred while processing the request.")
	}
}

func handleReceiptProcessingPanic(receiptUUID string) {
	if r := recover(); r != nil {
		database := data.GetDatabase()

		receipt := database.GetReceipt(receiptUUID)

		receipt.IsProcessing = false

		receipt.HasError = true

		receipt.Error = "Error extracting details from receipt."

		database.UpdateReceipt(receiptUUID, *receipt)

		panic("An error occurred while processing the receipt.")
	}
}

func validateReceipt(file *multipart.FileHeader, err error) {
	var response = errorResponse{}
	response.Code = http.StatusBadRequest
	if err != nil {
		response.Message = "File not found in the request."
		panic(response)
	}

	ext := filepath.Ext(file.Filename)
	if !isValidFileExtension(ext) {
		response.Message = "Invalid file format. Supported formats: PNG, JPG, PDF."
		panic(response)
	}
}

func saveFile(c *gin.Context, file *multipart.FileHeader, receiptUUID string) string {
	var response = errorResponse{}
	response.Code = http.StatusInternalServerError
	uploadDir := "./uploads/"
	fileName := receiptUUID + "_" + file.Filename
	// Save the uploaded file to the specified directory
	uploadPath := filepath.Join(uploadDir, fileName)
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		response.Message = "Error saving the file."
		panic(response)
	}
	return uploadPath
}

func processReceipt(receiptUUID string, filePath string) {
	defer handleReceiptProcessingPanic(receiptUUID)

	processor := receiptprocessor.GetReceiptProcessor()

	receiptDetails := processor.ExtractDetailsFromReceipt(filePath)

	if receiptDetails == "" {
		panic("Error extracting details from receipt.")
	}

	database := data.GetDatabase()

	receipt := database.GetReceipt(receiptUUID)

	receipt.IsProcessing = false

	receipt.Details = receiptDetails

	database.UpdateReceipt(receiptUUID, *receipt)

}
