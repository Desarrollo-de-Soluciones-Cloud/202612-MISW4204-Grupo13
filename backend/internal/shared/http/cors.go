package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const corsAllowHeaders = "Authorization, Content-Type, Accept"
const corsAllowMethods = "GET, POST, PUT, DELETE, OPTIONS"

func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	allowedOriginsMap := make(map[string]struct{}, len(allowedOrigins))
	allowAnyOrigin := len(allowedOrigins) == 0

	for _, origin := range allowedOrigins {
		allowedOriginsMap[origin] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && (allowAnyOrigin || isAllowedOrigin(origin, allowedOriginsMap)) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", corsAllowHeaders)
			c.Header("Access-Control-Allow-Methods", corsAllowMethods)
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isAllowedOrigin(origin string, allowedOrigins map[string]struct{}) bool {
	_, ok := allowedOrigins[origin]
	return ok
}
