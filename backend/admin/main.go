package main

import (
	"context"
	"data-collection/db"
	"os"

	"flag"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := godotenv.Load()
	if err != nil {
		log.Warn().Msg("couldn't loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URI"))
	if err != nil {
		log.Panic().Err(err).Msg("Failed to connect to database")
	} else {
		log.Info().Msg("Connected to database")
	}

	db_client := db.New(conn)

	var username = flag.String("username", "", "new admin username")
	var password = flag.String("password", "", "new admin password")

	flag.Parse()

	if *username == "" || *password == "" {
		log.Panic().Msg("Username and password are required")
	}

	password_hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to hash password")
	}

	id, err := db_client.InsertAdminUser(context.Background(), db.InsertAdminUserParams{
		Username:     *username,
		PasswordHash: string(password_hash),
	})
	if err != nil {
		log.Panic().Err(err).Msg("Failed to insert admin user")
	}

	log.Info().Int32("id", id).Str("username", *username).Msg("Admin user created successfully")
}
