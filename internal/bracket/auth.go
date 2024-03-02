package bracket

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func AzureAuthorizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientPrincipalID := os.Getenv("X-MS-CLIENT-PRINCIPAL-ID")
		clientPrincipalName := os.Getenv("X-MS-CLIENT-PRINCIPAL-Name")
		clientPrincipalIDP := os.Getenv("X-MS-CLIENT-PRINCIPAL-IDP")
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
