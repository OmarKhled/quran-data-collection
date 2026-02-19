package controllers

import (
	"context"
	"data-collection/db"

	"github.com/gin-gonic/gin"
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

func CreateUserTask(db_client *db.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		audio_file, err := c.FormFile("audio_file")
		if err != nil {
			log.Error().Err(err).Msg("error getting audio file")
			c.JSON(400, gin.H{
				"msg": "error getting audio file",
			})
			return
		}

		log.Info().Msgf("Received audio file: %s", audio_file.Filename)

		var body CreateUserTaskBody
		if err := c.ShouldBindJSON(&body); err != nil {
			log.Error().Err(err).Msg("error binding JSON")
			c.JSON(400, gin.H{
				"msg": err.Error(),
			})
			return
		}

		// err = db_client.CreateUserTask(context.Background(), db.CreateUserTaskParams{
		// 	UserID:        int32(body.UserID),
		// 	TaskID:        int32(body.TaskID),
		// 	AudioDuration: body.AudioDuration,
		// 	AudioFile:     audio_file.Filename,
		// })

	}
}
