package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type IndexController struct{}

var Index = IndexController{}

func (IndexController) index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "uptime": time.Now().UTC()})
}

func (IndexController) info(c *gin.Context) {
	uuid := c.Param("uuid")
	c.JSON(http.StatusOK, gin.H{"id": uuid, "name": "example"})
}
