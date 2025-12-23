package middleware

import (
	"github.com/gin-gonic/gin"
)

const CtxServerDataKey = "serverData"

func ServerDataMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		// gather/compute your server info
		s := map[string]interface{}{
			"Debug": gin.IsDebugging(),
			"Data": map[string]interface{}{
				"ClientIP": c.ClientIP(),
				"Method":   c.Request.Method,
				"Query":    c.Request.URL.Query(),
			},
			"Request": c.Request,
		}

		c.Set(CtxServerDataKey, s)
		c.Next()
	}
}
