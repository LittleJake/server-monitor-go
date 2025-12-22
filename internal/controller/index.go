package controller

import (
	"net/http"

	"github.com/LittleJake/server-monitor-go/internal/util"
	"github.com/gin-gonic/gin"
)

type IndexController struct{}

var Index = IndexController{}

func (IndexController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
	//c.JSON(http.StatusOK, gin.H{"status": "ok", "uptime": time.Now().UTC()})
}

func (IndexController) Info(c *gin.Context) {
	uuid := c.Param("uuid")

	info, err := util.GetInfo(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve info"})
		return
	}

	latest, err := util.GetCollectionLatest(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve latest collection data"})
		return
	}

	c.HTML(http.StatusOK, "info.html", gin.H{
		"uuid":     uuid,
		"base_url": util.GetEnv("BASE_URL", ""),
		"info":     info,
		"latest": latest,
	})
	// c.JSON(http.StatusOK, gin.H{"id": uuid, "name": "example"})
}
