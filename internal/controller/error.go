package controller

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

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
	errorData.(map[string]interface{})["Error"] = "404 Not Found"
	errorData.(map[string]interface{})["Message"] = fmt.Sprintf("%s %s", c.Request.URL.String(), "Not Found")
	c.HTML(http.StatusNotFound, "error.html", errorData)
}

func (ErrorController) NoMethodError(c *gin.Context) {
	errorData, ok := c.Get(middleware.CtxServerDataKey)
	if !ok {
		c.HTML(http.StatusMethodNotAllowed, "error.html", nil)
		return
	}
	errorData.(map[string]interface{})["Error"] = "405 Method Not Allowed"
	errorData.(map[string]interface{})["Message"] = fmt.Sprintf("%s %s", c.Request.URL.String(), "Method Not Allowed")
	c.HTML(http.StatusMethodNotAllowed, "error.html", errorData)
}

func (ErrorController) InternalServerError(c *gin.Context, recovered interface{}) {

	errorData, ok := c.Get(middleware.CtxServerDataKey)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	recoveredError, _ := recovered.(error)
	fmt.Println(recoveredError)

	// Generate stack trace
	buf := make([]byte, 40960)
	n := runtime.Stack(buf, true)
	stackTrace := string(buf[:n])
	stackTraceLines := strings.Split(stackTrace, "\n")

	errorData.(map[string]interface{})["Error"] = "500 Internal Server Error"
	errorData.(map[string]interface{})["Message"] = fmt.Sprintf("%s", recoveredError.Error())
	errorData.(map[string]interface{})["StackTraceLines"] = stackTraceLines
	c.HTML(http.StatusInternalServerError, "error.html", errorData)
}
