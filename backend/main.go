package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SuccessfulResponse(c *gin.Context, statusCode int, message string) {
	cookie, err := c.Cookie("auth")
	if err != nil {
		cookie = "random_stuff"
	}
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("auth", cookie, 3600, "/", "localhost", true, true)

	if message == "" {
		message = http.StatusText(statusCode)
	}

	c.JSON(statusCode, gin.H{
		"message": message,
	})
}

func main() {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true

	router.Use(sessions.Sessions("mysession", store))
	router.Use(cors.New(config))

	router.GET("/ping", func(c *gin.Context) {
		cookie, err := c.Cookie("auth")
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		SuccessfulResponse(c, 200, cookie)
	})

	router.GET("/get-cookie", func(c *gin.Context) {
		cookie, err := c.Cookie("auth")
		if err != nil {
			cookie = "random_stuff"
			c.SetSameSite(http.SameSiteNoneMode)
			c.SetCookie("auth", cookie, 3600, "/", "localhost", true, true)
		}
		c.JSON(200, gin.H{
			"message": "cookie got",
		})
	})
	log.Fatal(router.Run())
}
