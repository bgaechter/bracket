package bracket

import (
	"log"

	"github.com/gin-gonic/gin"
)

func AzureAuthorizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientPrincipalID := c.Request.Header.Get("X-Ms-Client-Principal-Id")
		clientPrincipalName := c.Request.Header.Get("X-Ms-Client-Principal-Name")
		clientPrincipalIDP := c.Request.Header.Get("X-Ms-Client-Principal-Idp")

		c.Next()

		// after request
		log.Println(clientPrincipalID)
		log.Println(clientPrincipalName)
		log.Println(clientPrincipalIDP)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}
