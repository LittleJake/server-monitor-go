package api

import (
	"github.com/gin-gonic/gin"
)

type ThermalAPI struct{}

var Thermal = ThermalAPI{}

func (ThermalAPI) get(c *gin.Context) {}
