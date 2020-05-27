package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {

	s, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("SECRETID"),  // id
			os.Getenv("SECRETKEY"), // secret
			""), // token can be left blank for now
	})

	// Create S3 service client
	svc := s3.New(s)

	// Get the list of items in bucket
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(os.Getenv("BUCKETNAME"))})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", os.Getenv("BUCKETNAME"), err)
	}

	for _, item := range resp.Contents {
		fmt.Println()
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
		fmt.Println()
	}

	fmt.Println("Found", len(resp.Contents), "items in bucket", os.Getenv("BUCKETNAME"))
	fmt.Println("")

}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
