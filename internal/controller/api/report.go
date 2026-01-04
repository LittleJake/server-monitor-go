package api

import (
	"github.com/gin-gonic/gin"
)

type ReportAPI struct{}

var Report = ReportAPI{}

func (ReportAPI) Set(c *gin.Context) {}
