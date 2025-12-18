package api

import (
	"github.com/gin-gonic/gin"
)

type MemoryAPI struct{}

var Memory = MemoryAPI{}

func (MemoryAPI) Get(c *gin.Context) {}
