package api

import (
	"net/http"
	"github.com/LittleJake/server-monitor-go/redis"
	"github.com/gin-gonic/gin"
)

type MemoryAPI struct{}

var Memory = MemoryAPI{}

func (MemoryAPI) Get(c *gin.Context) {
	data := redis.RedisGet(c.Request.Context(), redis.RedisClient, "memory")

	c.JSON(http.StatusOK, gin.H{
		"users": []gin.H{
			{"id": 1, "name": "alice"},
			{"id": 2, "name": "bob"},
		},
	})
}
