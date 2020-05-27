package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func Handlers() {
	// Default Router.

	// Middleware yang terpasang adalah Logger dan Recovery
	router := gin.Default()

	// Static
	router.Static("/", "./public")

	// Post Request Most Important is here
	router.POST("/uploadfiles", UploadFileToS3)

	// Run
	router.Run(":8080")
}

// UploadFileToS3 saves a file to aws bucket and returns the url to the file and an error if there's any
func UploadFileToS3(c *gin.Context) {

	// Multiple Form
	fmt.Println("-----------------------------------Uploading Files-----------------------------------------------------")
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("err: %s", err.Error()))
		return
	}

	// Files
	files := form.File["files"]
	if files == nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Could not get uploaded file"))
		return
	}

	// Session
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("SECRETID"),  // id
			os.Getenv("SECRETKEY"), // secret
			""), // token can be left blank for now
	})

	if err != nil {
		log.Println(err.Error())
		c.String(http.StatusBadRequest, fmt.Sprintf("Could not upload file"))
		return
	}

	// For range
	for _, file := range files {

		// fmt.Println(file.Filename)
		filename := file.Filename

		// get the file size and read
		// the file content into a buffer
		size := file.Size
		buffer := make([]byte, size)

		// uploading
		_, err := s3.New(s).PutObject(&s3.PutObjectInput{
			Bucket:               aws.String(os.Getenv("BUCKETNAME")),
			Key:                  aws.String(filename),
			ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
			Body:                 bytes.NewReader(buffer),
			ContentLength:        aws.Int64(int64(size)),
			ContentType:          aws.String(http.DetectContentType(buffer)),
			ContentDisposition:   aws.String("attachment"),
			ServerSideEncryption: aws.String("AES256"),
			StorageClass:         aws.String("INTELLIGENT_TIERING"),
		})

		if err != nil {
			log.Println(err.Error())
			c.String(http.StatusBadRequest, fmt.Sprintf("Could not upload file"))
			return
		}
		url := "https://%s.s3-%s.amazonaws.com/%s"
		url = fmt.Sprintf(url, os.Getenv("BUCKETNAME"), os.Getenv("REGION"), filename)
		c.String(http.StatusBadRequest, fmt.Sprintf("Image uploaded successfully: %v", filename))
		fmt.Printf("Uploaded File Url %s\n", url)
	}

	// Response
	c.String(http.StatusOK, fmt.Sprintf("Files count : %d", len(files)))
	fmt.Println("----------------------------------------------------------------------------------------")
}

func main() {
	Handlers()
}
