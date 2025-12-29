package api

import (
	"net/http"

	"github.com/LittleJake/server-monitor-go/internal/util"
	"github.com/gin-gonic/gin"
)

type BatteryAPI struct{}

var Battery = BatteryAPI{}

func (BatteryAPI) Get(c *gin.Context) {
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
	battery := util.CollectionFormat(result, "Battery")
	c.JSON(http.StatusOK, battery)

}
