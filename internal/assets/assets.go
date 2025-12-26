package assets

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed favicon.ico
var faviconFS embed.FS

//go:embed locales
var LocalesFS embed.FS

//go:embed templates
var TemplatesFS embed.FS

// ServeFavicon serves the embedded favicon.ico
func ServeFavicon(c *gin.Context) {
	data, err := faviconFS.ReadFile("favicon.ico")
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Header("Content-Type", "image/x-icon")
	c.Data(http.StatusOK, "image/x-icon", data)
}
