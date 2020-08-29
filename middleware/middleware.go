package middleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
	_ "github.com/lib/pq"
)

func AuthMiddleware(c *gin.Context) {
	log.Println("start authorization middleware")
	authKey := c.GetHeader("Authorization")
	if authKey != "November 10, 2009" {
		c.JSON(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		c.Abort()
		return
	}
	c.Next()
	log.Println("end middleware")
}