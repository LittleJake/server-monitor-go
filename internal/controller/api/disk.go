package api

import (
	"github.com/gin-gonic/gin"
)

type DiskAPI struct{}

var Disk = DiskAPI{}

func (DiskAPI) Get(c *gin.Context) {}
