package assets

import (
	"embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed static/**
var StaticFS embed.FS

//go:embed favicon.ico
var faviconFS embed.FS

//go:embed manifest.json
var manifestFS embed.FS

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

func ServeManifest(c *gin.Context) {
	data, err := manifestFS.ReadFile("manifest.json")
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Header("Content-Type", "html/json")
	c.Data(http.StatusOK, "html/json", data)
}
