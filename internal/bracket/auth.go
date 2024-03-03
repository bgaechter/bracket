package bracket

import (
	"log"
<<<<<<< HEAD
=======
	"os"
>>>>>>> f131f51aa6430330956a6a25ee7ca830619eb296

	"github.com/gin-gonic/gin"
)

func AzureAuthorizer() gin.HandlerFunc {
	return func(c *gin.Context) {
<<<<<<< HEAD
		clientPrincipalID := c.Request.Header.Get("X-Ms-Client-Principal-Id")
		clientPrincipalName := c.Request.Header.Get("X-Ms-Client-Principal-Name")
		clientPrincipalIDP := c.Request.Header.Get("X-Ms-Client-Principal-Idp")

=======
		clientPrincipalID := os.Getenv("X-MS-CLIENT-PRINCIPAL-ID")
		clientPrincipalName := os.Getenv("X-MS-CLIENT-PRINCIPAL-Name")
		clientPrincipalIDP := os.Getenv("X-MS-CLIENT-PRINCIPAL-IDP")
>>>>>>> f131f51aa6430330956a6a25ee7ca830619eb296
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
