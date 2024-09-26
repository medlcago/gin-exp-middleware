package main

import (
	"gin-exp-middleware/exp"
	"gin-exp-middleware/middleware"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type User struct {
	Username string `json:"username" binding:"required,ascii,min=4,max=16"`
	Email    string `json:"email" binding:"required,email"`
}

func main() {
	router := gin.Default()
	router.Use(middleware.ExpMiddleware())

	router.POST("/", func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")
		if apiKey != "test123" {
			c.Error(exp.NewHttpExp(http.StatusUnauthorized, "invalid API key"))
			return
		}

		user := new(User)
		if err := c.ShouldBindJSON(&user); err != nil {
			c.Error(exp.NewValidationExp(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"user":    user,
		})
	})

	log.Println("Server started on :8080")
	log.Fatal(router.Run(":8080"))
}
