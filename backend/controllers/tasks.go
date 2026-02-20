package controllers

import (
	"context"
	"data-collection/db"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func GetTask(db_client *db.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		task, err := db_client.GetTask(context.Background())

		if err != nil {
			log.Error().Err(err).Msg("error fetching task")
			c.JSON(500, gin.H{
				"msg": "error fetching task",
			})
			return
		}

		c.JSON(200, gin.H{
			"task": task,
		})
	}
}

type CreateUserTaskBody struct {
	UserID        int     `json:"user_id" binding:"required"`
	TaskID        int     `json:"task_id" binding:"required"`
	AudioDuration float64 `json:"audio_duration" binding:"required"`
}

func CreateUserTask(db_client *db.Queries, r2_client *s3.Client, bucketName string) func(c *gin.Context) {
	return func(c *gin.Context) {
		audio_file, err := c.FormFile("audio")
		if err != nil {
			log.Error().Err(err).Msg("error getting audio file")
			c.JSON(400, gin.H{
				"msg": "error getting audio file",
			})
			return
		}

		var body CreateUserTaskBody
		userId := c.PostForm("user_id")
		taskId := c.PostForm("task_id")
		audioDuration := c.PostForm("duration")

		if userId == "" || taskId == "" || audioDuration == "" {
			log.Error().Msg("missing required fields in form data")
			c.JSON(400, gin.H{
				"msg": "missing required fields in form data",
			})
			return
		}

		body.UserID, _ = strconv.Atoi(userId)
		body.TaskID, _ = strconv.Atoi(taskId)
		body.AudioDuration, _ = strconv.ParseFloat(audioDuration, 64)

		task, err := db_client.GetTaskByID(context.Background(), int32(body.TaskID))

		if err != nil {
			log.Error().Err(err).Msg("error fetching task")
			c.JSON(500, gin.H{
				"msg": "error creating task",
			})
			return
		}

		uuid := uuid.New()
		audio_extension := filepath.Ext(audio_file.Filename)
		// filename: surah_ayah_userid_uuid
		filename := fmt.Sprintf("%d_%d_%d_%s%s", task.Surah, task.Ayah, body.UserID, uuid.String(), audio_extension)

		file, err := audio_file.Open()
		if err != nil {
			log.Error().Err(err).Msg("error opening audio file")
			c.JSON(500, gin.H{
				"msg": "error opening audio file",
			})
			return
		}
		defer file.Close()

		contentType := audio_file.Header.Get("Content-Type")
		log.Info().Msgf("uploading file %s to R2 bucket %s, content-type: %s", filename, bucketName, contentType)
		out, err := r2_client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket:        &bucketName,
			Key:           &filename,
			Body:          file,
			ContentType:   &contentType,
			ContentLength: &audio_file.Size,
		})

		if err != nil {
			log.Error().Err(err).Msg("error uploading audio file")
			c.JSON(500, gin.H{
				"msg": "error uploading audio file",
			})
			return
		}

		log.Info().Msgf("uploaded file to R2 with ETag %s", *out.ETag)

		err = db_client.CreateUserTask(context.Background(), db.CreateUserTaskParams{
			UserID:        int32(body.UserID),
			TaskID:        int32(body.TaskID),
			AudioDuration: body.AudioDuration,
			AudioUrl:      bucketName + "/" + filename,
		})

		c.JSON(200, gin.H{
			"msg": "success",
		})
	}
}
