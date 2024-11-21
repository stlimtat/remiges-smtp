package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleAuth(
	c *gin.Context,
) {
	if c.Request.Header.Get("Authorization") != "foobar" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	c.Next()
}
