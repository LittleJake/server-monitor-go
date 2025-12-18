package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LittleJake/server-monitor-go/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type MemoryAPI struct{}

var Memory = MemoryAPI{}

func (MemoryAPI) Get(c *gin.Context) {
	uuid, ok := c.Params.Get("uuid")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "uuid parameter is required",
		})
		return
	}

	data, err := util.RedisZRangeByScoreWithScores(
		c.Request.Context(),
		util.RedisClient,
		"system_monitor:collection:"+uuid,
		&redis.ZRangeBy{Min: "0", Max: fmt.Sprint(time.Now().Unix())},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"memory": data,
	})
}
