package controllers

import (
	"data-collection/db"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(db_client *db.Queries) func(c *gin.Context) {
	secret := os.Getenv("SECRET")
	frontend_domain := os.Getenv("FRONTEND_ROOT_DOMAIN")
	return func(c *gin.Context) {
		var body LoginBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{
				"msg": "invalid request body",
			})
			return
		}

		user, err := db_client.GetAdminUser(c.Request.Context(), body.Username)
		if err != nil {
			c.JSON(401, gin.H{
				"msg": "invalid username or password",
			})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))
		if err != nil {
			c.JSON(401, gin.H{
				"msg": "invalid username or password",
			})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":       user.ID,
			"username": user.Username,
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			c.JSON(500, gin.H{
				"msg": "error generating token",
			})
			return
		}

		var domain string
		if os.Getenv("ENV") == "production" {
			domain = fmt.Sprintf(".%s", frontend_domain)
		} else {
			domain = ""
		}

		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie("token", tokenString, 3600*24, "/", domain, true, true)

		c.JSON(200, gin.H{
			"msg": "login successful",
		})
	}
}

func Logout(db_client *db.Queries) func(c *gin.Context) {
	frontend_domain := os.Getenv("FRONTEND_ROOT_DOMAIN")
	return func(c *gin.Context) {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "token",
			Value:    "",
			Path:     "/",
			Domain:   fmt.Sprintf(".%s", frontend_domain),
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		})

		c.JSON(http.StatusOK, map[string]string{
			"msg": "logout successful",
		})
	}
}
