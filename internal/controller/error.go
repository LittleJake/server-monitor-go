package controller

import (
	"fmt"
	"net/http"

	"github.com/LittleJake/server-monitor-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

type ErrorController struct{}

var Error = ErrorController{}

func (ErrorController) NoRouteError(c *gin.Context) {
	errorData, ok := c.Get(middleware.CtxServerDataKey)
	if !ok {
		c.HTML(http.StatusNotFound, "error.html", nil)
		return
	}
	errorData.(map[string]interface{})["Message"] = fmt.Sprintf("%s %s", c.Request.URL.String(), "Not Found")
	c.HTML(http.StatusNotFound, "error.html", errorData)
}

func (ErrorController) NoMethodError(c *gin.Context) {
	errorData, ok := c.Get(middleware.CtxServerDataKey)
	if !ok {
		c.HTML(http.StatusMethodNotAllowed, "error.html", nil)
		return
	}
	errorData.(map[string]interface{})["Message"] = fmt.Sprintf("%s %s", c.Request.URL.String(), "Method Not Allowed")
	c.HTML(http.StatusMethodNotAllowed, "error.html", errorData)
}
