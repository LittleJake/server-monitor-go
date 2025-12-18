package middleware

import (
	"github.com/LittleJake/server-monitor-go/internal/util"
	"github.com/gin-gonic/gin"
)

const CtxServerDataKey = "serverData"

func ServerDataMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		// gather/compute your server info
		s := map[string]interface{}{
			"Debug": gin.IsDebugging(),
			"Envs":  util.GetAllEnvs(),
		}

		c.Set(CtxServerDataKey, s)
		c.Next()
	}
}
