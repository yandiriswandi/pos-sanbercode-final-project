package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yandiriswandi/pos-sanbercode-final-project/models"
)

// Load .env variables
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func UploadFile(c *gin.Context) {
	// Get the file from the request
	file, _ := c.FormFile("file")
	timestamp := time.Now().Unix()

	// Initialize Cloudinary
	cloudinaryURL := os.Getenv("CLOUDINARY_URL")
	cld, err := cloudinary.NewFromURL(cloudinaryURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Code: 500, Status: "failed", Message: "Response to connect to Cloudinary"})
		return
	}

	// Upload the file
	folder := "sanbercode"
	publicID := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	uploadResult, err := cld.Upload.Upload(c, file, uploader.UploadParams{
		Folder:       folder,
		PublicID:     publicID,
		ResourceType: "image",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Code: 400, Status: "failed", Message: "failed to upload image"})
		return
	}

	// Return the URL of the uploaded image
	c.JSON(http.StatusOK, models.SuccessAddUpdate{Data: uploadResult.URL, Status: "success", Code: 200, Message: "success upload file"})
}
