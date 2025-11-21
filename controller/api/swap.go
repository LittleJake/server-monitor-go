package api

import (
	"github.com/gin-gonic/gin"
)

type SwapAPI struct{}

var Swap = SwapAPI{}

func (SwapAPI) get(c *gin.Context) {}
