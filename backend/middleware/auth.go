package middleware

import (
	"data-collection/db"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func IsAuthorized(db_client *db.Queries) func(c *gin.Context) {
	secret := os.Getenv("SECRET")
	return func(c *gin.Context) {
		token_cookie, err := c.Cookie("token")

		if err == http.ErrNoCookie {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "authorization token not provided",
			})
			c.Abort()
			return
		} else if err != nil {
			log.Error().Err(err).Msg("error retrieving authorization cookie")
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "error processing authorization token",
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(token_cookie, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		}, jwt.WithJSONNumber())

		if err != nil {
			log.Error().Err(err).Msg("error parsing JWT token")
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "invalid token",
			})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "invalid or expired token",
			})
			c.Abort()
			return
		}

		payload, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Error().Msg("invalid token claims type")
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "error processing authorization token",
			})
			c.Abort()
			return
		}

		user_id, err := payload["id"].(json.Number).Int64()
		if err != nil {
			log.Error().Err(err).Msg("couldn't parse user id from token")
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "error processing authorization token",
			})
			c.Abort()
			return
		}

		username, ok := payload["username"].(string)
		if !ok {
			log.Error().Msg("Couldn't parse username from token")
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "error processing authorization token",
			})
			c.Abort()
			return
		}

		user, err := db_client.GetAdminUser(c.Request.Context(), username)
		if err == pgx.ErrNoRows {
			log.Info().Msgf("User does not exist for token username: %v", username)
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "unauthorized access",
			})
			c.Abort()
			return
		} else if err != nil {
			log.Error().Err(err).Msg("error fetching user from db")
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "error processing authorization token",
			})
			c.Abort()
			return
		}
		if int64(user.ID) != user_id || user.Username != username {
			log.Info().Msgf("User ID or Username mismatch: expected ID %v, got %v; expected username %v, got %v", user_id, user.ID, username, user.Username)
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "unauthorized access",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
