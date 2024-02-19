package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bgaechter/bracket/internal/bracket"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		principalName := c.Request.Header.Get("X-MS-CLIENT-PRINCIPAL-NAME")
		c.Next()
		log.Println(principalName)
		authorizedUsers := os.Getenv("BRACKET_USERS")
		if strings.Contains(authorizedUsers, principalName) {
			c.Set("admin", "true")
		} else {
			c.Redirect(http.StatusTemporaryRedirect, "/unauthorized")
		}
	}
}

func main() {

	router := gin.Default()
	store := memstore.NewStore([]byte("ranker"))
	router.Use(sessions.Sessions("session", store))
	// router.Use(Auth())

	router.LoadHTMLGlob("templates/*")
	router.GET("/", bracket.GetIndex)
	router.GET("/unauthorized", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "unauthorized"})
	})
	router.POST("/", func(c *gin.Context) {
		bracket.PostPlay(c)
		c.Request.URL.Path = "/"
		c.Request.Method = "GET"
		router.HandleContext(c)
	})
	router.POST("/postMatch", func(c *gin.Context) {
		bracket.PostMatch(c)
		c.Request.URL.Path = "/"
		c.Request.Method = "GET"
		router.HandleContext(c)
	})

	router.Run("localhost:8080")
}
