package api

import (
	"github.com/gin-gonic/gin"
)

type ThermalAPI struct{}

var Thermal = ThermalAPI{}

func (ThermalAPI) Get(c *gin.Context) {}
