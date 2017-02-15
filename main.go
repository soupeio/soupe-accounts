package main

import (
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

type Register struct {
	UserName             string `form:"username" binding:"required"`
	Email                string `form:"email"    binding:"required"`
	Password             string `form:"password" binding:"required"`
	PasswordConfirmation string `form:"password_confirmation" binding:"required"`
}

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.tmpl", gin.H{})
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
	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
