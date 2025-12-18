package api

import (
	"github.com/gin-gonic/gin"
)

type CpuAPI struct{}

var Cpu = CpuAPI{}

func (CpuAPI) Get(c *gin.Context) {}
