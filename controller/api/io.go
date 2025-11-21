package api

import (
	"github.com/gin-gonic/gin"
)

type IOAPI struct{}

var IO = IOAPI{}

func (IOAPI) get(c *gin.Context) {}
