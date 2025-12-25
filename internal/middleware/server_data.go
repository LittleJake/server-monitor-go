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
			"Context": c,
		}

		// set lang cookie
		lang := c.Query("lang")
		if lang == "" {
			if cookie, err := c.Cookie("lang"); err == nil {
				lang = cookie
			}
		}
		c.SetCookie("lang", lang, 0, "/", c.Request.URL.Hostname(), false, true)

		c.Set(CtxServerDataKey, s)
		c.Next()
	}
}
