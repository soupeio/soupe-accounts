package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

type Register struct {
	UserName             string `form:"username" binding:"required"`
	Email                string `form:"email"    binding:"required"`
	Password             string `form:"password" binding:"required"`
	PasswordConfirmation string `form:"password_confirmation" binding:"required"`
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.Static("/assets", "./assets")

	router.GET("/register", func(c *gin.Context) {
		cookie, _ := c.Cookie("id")
		c.HTML(http.StatusOK, "register.tmpl", gin.H{"cookie": cookie})
	})

	router.POST("/register", func(c *gin.Context) {

		var form Register
		if c.Bind(&form) == nil {
			c.HTML(http.StatusOK, "register_success.tmpl", gin.H{
				"username": form.UserName,
				"email":    form.Email,
				"password": form.Password,
			})
		}

	})

	router.GET("/get_token", func(c *gin.Context) {
		random_bytes, _ := GenerateRandomString(128)

		c.SetCookie("id", random_bytes, 3600, "/", "", false, false)

		c.HTML(http.StatusOK, "get_token.tmpl", gin.H{
			"token_string": "Hallo",
			"token":        "Ballo",
			"claims":       random_bytes,
		})

	})

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
