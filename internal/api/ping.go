package api

import (
	"github.com/gin-gonic/gin"
)

// PingHandler обработчик GET /api/ping
func PingHandler(c *gin.Context) {
	c.String(200, "ok")
}
