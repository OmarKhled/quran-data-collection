package main

import (
	"context"
	"data-collection/config"
	"data-collection/db"
	"data-collection/routes"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := godotenv.Load()
	if err != nil {
		log.Warn().Msg("couldn't loading .env file")
	}

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "healthy",
		})
	})

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URI"))
	if err != nil {
		log.Panic().Err(err).Msg("Failed to connect to database")
	} else {
		log.Info().Msg("Connected to database")
	}

	db_client := db.New(conn)
	r2_client := config.GetR2Client()
	log.Info().Msg("Initialized R2 client")

	routes.SetupRoutes(router, db_client, r2_client)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}

	log.Info().Msgf("Starting server on port %s", PORT)
	router.Run(":" + PORT)
}
