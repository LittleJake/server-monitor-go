package controller

import (
	"net/http"

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
	c.HTML(http.StatusOK, "info.html", gin.H{
		"uuid": uuid,
	})
	// c.JSON(http.StatusOK, gin.H{"id": uuid, "name": "example"})
}
