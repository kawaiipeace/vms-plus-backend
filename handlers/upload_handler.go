package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"vms_plus_be/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct {
}

func GenerateFileName() string {
	// Generate a unique ID using UUID
	uniqueID := uuid.New().String()

	// Use the current timestamp for additional uniqueness
	timestamp := time.Now().Unix()

	// Combine the timestamp and UUID to create a unique file name
	return fmt.Sprintf("%d_%s", timestamp, uniqueID)
}

// @Summary Uploads a file
// @Description Upload a file to the server and save it in D:\uploads
// @Tags Uploads
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]interface{} "Successfully uploaded file"
// @Failure 400 {object} map[string]interface{} "Invalid file"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/upload [post]
func (h *UploadHandler) UploadFile(c *gin.Context) {
	// Maximum file size limit (e.g., 10MB)
	const maxFileSize = 10 << 20 // 10MB

	// Set the maximum allowed size for uploads
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxFileSize)

	// Get the file from the form
	file, _ := c.FormFile("file")
	if file == nil {
		c.JSON(400, gin.H{
			"error": "No file is provided",
		})
		return
	}

	// Generate a new file name
	fileName := GenerateFileName()
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		// If no extension is found, default to .bin
		ext = ".bin"
	}

	// Append the extension to the generated file name
	fileNameWithExt := fileName + ext

	uploadDir := config.AppConfig.SaveFilePath

	// Create the uploads directory if it doesn't exist
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Unable to create upload directory",
		})
		return
	}

	// Save the file to D:\uploads
	filePath := filepath.Join(uploadDir, fileNameWithExt)
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Unable to save file",
		})
		return
	}

	// Send success response with file name
	c.JSON(200, gin.H{
		"message":  "File uploaded successfully",
		"file_url": config.AppConfig.SaveFileUrl + "/" + fileNameWithExt,
	})
}
