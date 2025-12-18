package api

import (
	"github.com/gin-gonic/gin"
)

type BatteryAPI struct{}

var Battery = BatteryAPI{}

func (BatteryAPI) Get(c *gin.Context) {}
