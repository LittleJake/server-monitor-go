package api

import (
	"github.com/gin-gonic/gin"
)

type NetworkAPI struct{}

var Network = NetworkAPI{}

func (NetworkAPI) Get(c *gin.Context) {}
