package api

import (
	"github.com/gin-gonic/gin"
)

type ReportAPI struct{}

var Report = ReportAPI{}

func (ReportAPI) get(c *gin.Context) {}
