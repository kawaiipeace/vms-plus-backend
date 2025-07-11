package handlers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"vms_plus_be/config"
	"vms_plus_be/messages"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type UploadHandler struct {
}

var minioClient *minio.Client

// Initialize MinIO client
func InitMinIO(endpoint, accessKey, secretKey string, useSSL bool) {
	if config.AppConfig.MinIoEndPoint == "" {
		return
	}
	var err error
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}
	log.Println("MinIO client initialized successfully")
}
func GenerateFileName(fileExt string) string {
	// Generate a unique ID using UUID
	uniqueID := uuid.New().String()

	// Use the current timestamp for additional uniqueness
	timestamp := time.Now().Unix()

	// Combine the timestamp and UUID to create a unique file name
	return fmt.Sprintf("%d_%s%s", timestamp, uniqueID, fileExt)
}

// @Summary Uploads a file
// @Description Upload a file to the server and save it in D:\uploads
// @Tags Uploads
// @Accept multipart/form-data
// @Produce json
// @Security ApiKeyAuth
// @Param file formData file true "File to upload"
// @Router /api/upload [post]
func (h *UploadHandler) UploadFile(c *gin.Context) {
	if config.AppConfig.MinIoEndPoint == "" {
		DevUploadFile(c)
		return
	}

	bucketName := config.AppConfig.MinIoAccessKey

	// Parse the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required", "message": messages.ErrInvalidJSONInput.Error()})
		return
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file", "message": messages.ErrInternalServer.Error()})
		return
	}
	defer src.Close()

	// Create bucket if it doesn't exist
	ctx := context.Background()

	// Upload the file to MinIO
	fileExt := filepath.Ext(file.Filename)
	fileName := GenerateFileName(fileExt)
	_, err = minioClient.PutObject(ctx, bucketName, fileName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error(), "bucket_name": bucketName, "message": messages.ErrInternalServer.Error()})
		return
	}

	// Generate API URL for viewing the file
	apiHost := c.Request.Host // Get the host address from the incoming request
	fileURL := "https://" + apiHost + "/api/files/" + bucketName + "/" + fileName

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message":   "File uploaded successfully",
		"file_name": file.Filename,
		"file_url":  fileURL,
	})
}

func determineContentType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".svg":
		return "image/svg+xml"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream" // Default binary stream
	}
}

func (h *UploadHandler) GetFile(c *gin.Context) {
	// Extract bucket name and file name from the route parameters
	bucketName := c.Param("bucket")
	fileName := c.Param("file")

	// Get the object from MinIO
	object, err := minioClient.GetObject(context.Background(), bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file", "message": messages.ErrInternalServer.Error()})
		return
	}

	// Set the appropriate Content-Type
	contentType := determineContentType(fileName)
	c.Header("Content-Type", contentType)

	// Serve the object content as a response
	_, err = io.Copy(c.Writer, object)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content", "message": messages.ErrInternalServer.Error()})
	}
}

func (h *UploadHandler) ListFiles(c *gin.Context) {
	// Extract bucket name from route parameters
	bucketName := c.Param("bucket")
	fmt.Print(bucketName)

	// Channel for receiving object information
	objectCh := minioClient.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{})

	// Prepare a slice to hold object information
	var files []map[string]interface{}

	// Iterate through the objects in the bucket
	for object := range objectCh {
		if object.Err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files", "details": object.Err.Error(), "message": messages.ErrInternalServer.Error()})
			return
		}
		files = append(files, map[string]interface{}{
			"name":         object.Key,
			"size":         object.Size,
			"lastModified": object.LastModified,
		})
	}

	// Respond with the list of files
	c.JSON(http.StatusOK, gin.H{
		"bucket": bucketName,
		"files":  files,
	})
}

func DevUploadFile(c *gin.Context) {
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
	ext := filepath.Ext(file.Filename)

	if ext == "" {
		// If no extension is found, default to .bin
		ext = ".bin"
	}
	fileNameWithExt := GenerateFileName(ext)

	uploadDir := config.AppConfig.DevSaveFilePath
	// Create the uploads directory if it doesn't exist
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "Unable to create upload directory",
			"message": messages.ErrInternalServer.Error(),
		})
		return
	}
	fmt.Println(uploadDir)
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
		"message":   "File uploaded successfully",
		"file_name": file.Filename,
		"file_url":  config.AppConfig.DevSaveFileUrl + "/" + fileNameWithExt,
	})
}
