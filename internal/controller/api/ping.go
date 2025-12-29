package api

import (
	"net/http"

	"github.com/LittleJake/server-monitor-go/internal/util"
	"github.com/gin-gonic/gin"
)

type PingAPI struct{}

var Ping = PingAPI{}

func (PingAPI) Get(c *gin.Context) {
	uuid, ok := c.Params.Get("uuid")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "uuid parameter is required",
		})
		return
	}

	result, err := util.GetCollection(uuid, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ping := util.CollectionFormat(result, "Ping")
	c.JSON(http.StatusOK, ping)
}
