package api

import (
	"github.com/gin-gonic/gin"
)

type PingAPI struct{}

var Ping = PingAPI{}

func (PingAPI) get(c *gin.Context) {}
