package controllers

import (
	"data-collection/db"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetUser(db_client *db.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		user_email := c.Query("email")
		user_id, err := db_client.GetUser(c, user_email)
		if err != nil {
			c.JSON(400, gin.H{
				"msg": "User not found",
			})
			return
		}

		c.JSON(200, gin.H{
			"id": user_id,
		})
	}
}

type CreateUserRequest struct {
	Name               string    `json:"name" binding:"required"`
	Email              string    `json:"email" binding:"required,email"`
	Country            string    `json:"country" binding:"required"`
	Province           string    `json:"province" binding:"required"`
	Age                int32     `json:"age" binding:"required"`
	Gender             db.Gender `json:"gender" binding:"required"`
	ProficiencyLevel   string    `json:"proficiency_level" binding:"required"`
	StudiedQuranBefore *bool     `json:"studied_quran_before" binding:"required"`
	JobTitle           string    `json:"job_title"`
	Lisper             string    `json:"lisper"`
}

func CreateUser(db_client *db.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		var body CreateUserRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{
				"msg": err.Error(),
			})
			return
		}

		user_id, err := db_client.CreateUser(c, db.CreateUserParams{
			Name:               body.Name,
			Email:              body.Email,
			Country:            body.Country,
			Province:           body.Province,
			Age:                body.Age,
			Gender:             body.Gender,
			ProficiencyLevel:   body.ProficiencyLevel,
			StudiedQuranBefore: *body.StudiedQuranBefore,
			JobTitle:           pgtype.Text{String: body.JobTitle, Valid: body.JobTitle != ""},
			Lisper:             pgtype.Text{String: body.Lisper, Valid: body.Lisper != ""},
		})

		if err != nil {
			c.JSON(500, gin.H{
				"msg": fmt.Sprintf("failed to create user: %v", err),
			})
			return
		}

		c.JSON(201, gin.H{
			"msg": "user created successfully",
			"id":  user_id,
		})
	}
}
