package main

import (
	"github.com/bgaechter/bracket/internal/bracket"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	store := memstore.NewStore([]byte("ranker"))
	router.Use(sessions.Sessions("session", store))

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

	router.Run("0.0.0.0:8080")
}
