package controllers

import (
	"data-collection/db"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

func getS3PathParts(s3Path string) (bucket, key string) {
	// Assuming the s3Path is in the format "bucket/key"
	parts := strings.SplitN(s3Path, "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func GetUsersTasks(db_client *db.Queries, r2_client *s3.Client) func(c *gin.Context) {
	return func(c *gin.Context) {
		tasks, err := db_client.GetUsersTasks(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{
				"msg": "error fetching user tasks",
			})
			return
		}

		presignClient := s3.NewPresignClient(r2_client)

		for i, task := range tasks {
			bucket, key := getS3PathParts(task.AudioUrl)
			if bucket == "" || key == "" {
				c.JSON(500, gin.H{
					"msg": "can't serve tasks",
				})
				return
			}

			expires := time.Now().Add(10 * time.Minute)

			presignResult, err := presignClient.PresignGetObject(c.Request.Context(), &s3.GetObjectInput{
				Bucket:          &bucket,
				Key:             &key,
				ResponseExpires: &expires,
			})
			if err != nil {
				c.JSON(500, gin.H{
					"msg": "error generating presigned URL",
				})
				return
			}
			tasks[i].AudioUrl = presignResult.URL
		}

		c.JSON(200, gin.H{
			"tasks": tasks,
		})
	}
}
