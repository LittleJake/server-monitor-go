package api

import (
	"github.com/gin-gonic/gin"
)

type MemoryAPI struct{}

var Memory = MemoryAPI{}

func (MemoryAPI) get(c *gin.Context) {}
