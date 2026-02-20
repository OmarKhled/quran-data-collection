package controllers

import (
	"data-collection/db"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
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

type UserPerformance struct {
	Name          string  `json:"name"`
	TotalDuration float64 `json:"total_duration"`
	Email         string  `json:"email"`
	Rank          int     `json:"rank"`
}

func GetUsersRanks(db_client *db.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		user_id := c.Query("user_id")
		ranks, err := db_client.GetUsersRanks(c.Request.Context())
		if err != nil {
			log.Error().Err(err).Msg("failed to get users ranks")
			c.JSON(500, gin.H{
				"msg": "failed to get users ranks",
			})
			return
		}

		if user_id != "" {
			userID, _ := strconv.Atoi(user_id)
			already_in_top_10 := false
			for _, p := range ranks {
				if p.ID == int32(userID) {
					already_in_top_10 = true
					break
				}
			}

			if !already_in_top_10 {
				user_rank, err := db_client.GetUserRank(c.Request.Context(), int32(userID))
				if err != nil {
					log.Error().Err(err).Msg("failed to get user rank")
				} else {
					ranks = append(ranks, db.GetUsersRanksRow(user_rank))
				}
			}
		}

		users_ranks := make([]UserPerformance, len(ranks))
		for i, p := range ranks {
			total_duration, _ := strconv.ParseFloat(fmt.Sprintf("%v", p.TotalDuration), 64)
			users_ranks[i] = UserPerformance{
				Name:          p.Name,
				TotalDuration: total_duration,
				Email:         p.Email,
				Rank:          int(p.Rank),
			}
		}

		c.JSON(200, gin.H{
			"users": users_ranks,
		})
	}
}

func GetUserRank(db_client *db.Queries) func(c *gin.Context) {
	return func(c *gin.Context) {
		user_id := c.Query("user_id")
		if user_id == "" {
			c.JSON(400, gin.H{
				"msg": "user_id is required",
			})
			return
		}

		userID, _ := strconv.Atoi(user_id)
		rank, err := db_client.GetUserRank(c.Request.Context(), int32(userID))
		if err != nil {
			log.Error().Err(err).Msg("failed to get user rank")
			c.JSON(500, gin.H{
				"msg": "failed to get user rank",
			})
			return
		}

		total_duration, err := db_client.GetUserTotalDuration(c.Request.Context(), int32(userID))
		if err != nil {
			log.Error().Err(err).Msg("failed to get user total duration")
			c.JSON(500, gin.H{
				"msg": "failed to get user rank",
			})
			return
		}

		c.JSON(200, gin.H{
			"user":           rank,
			"total_duration": total_duration,
		})
	}
}
