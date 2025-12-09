package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/LittleJake/server-monitor-go/controller"
	"github.com/LittleJake/server-monitor-go/controller/api"

	"github.com/gin-gonic/gin"
)

// SetupRouter builds and returns a gin.Engine with example routes and middleware.
func SetupRouter() *gin.Engine {
	// parse bool from environment variable `WEB.IS_DEBUG`.
	// Accepts common boolean string values like: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
	isDebug := false
	if b := GetEnvBool("IS_DEBUG", false); b {
		isDebug = b
	}

	println("Web Debug Mode:", isDebug)
	if isDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Built-in middleware
	r.Use(gin.Logger())
	// Use the default Recovery in debug mode so developers see panics.
	// In non-debug (production) mode, use a custom recovery that renders
	// a friendly error page instead of exposing stack traces.
	if gin.IsDebugging() {
		r.Use(gin.Recovery())
	} else {
		r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{})
		}))
	}

	// In non-debug mode, route unknown paths and methods to the error page.
	if !gin.IsDebugging() {
		r.NoRoute(func(c *gin.Context) {
			c.HTML(http.StatusNotFound, "error.html", gin.H{})
		})
		r.NoMethod(func(c *gin.Context) {
			c.HTML(http.StatusMethodNotAllowed, "error.html", gin.H{})
		})
	} else {
		data := map[string]interface{}{
			"Envs": GetAllEnvs(),
		}
		r.NoRoute(func(c *gin.Context) {
			data["Message"] = fmt.Sprintf("%s %s", c.Request.URL.String(), "Not Found")
			c.HTML(http.StatusNotFound, "error.html", data)
		})
		r.NoMethod(func(c *gin.Context) {
			data["Message"] = fmt.Sprintf("%s %s", c.Request.URL.String(), "Method Not Allowed")
			c.HTML(http.StatusMethodNotAllowed, "error.html", data)
		})
	}

	// Simple CORS middleware (for example purposes)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.GET("/", controller.Index.Index)

	// Info routes
	r.GET("/info/:uuid", controller.Index.Info)

	// API group
	_api := r.Group("/api")
	{
		_api.GET("/cpu/:uuid", api.Cpu.Get)
		_api.GET("/memory/:uuid", api.Memory.Get)
		_api.GET("/disk/:uuid", api.Disk.Get)
		_api.GET("/network/:uuid", api.Network.Get)
		_api.GET("/io/:uuid", api.IO.Get)
		_api.GET("/ping/:uuid", api.Ping.Get)
		_api.GET("/swap/:uuid", api.Swap.Get)
		_api.GET("/thermal/:uuid", api.Thermal.Get)
		_api.GET("/report/:uuid", api.Report.Get)
	}

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "uptime": time.Now().UTC()})
	})
	// Public routes
	r.GET("/404", func(c *gin.Context) {
		c.HTML(http.StatusNotFound,
			"error.html", gin.H{})
	})
	r.GET("/metrics", func(c *gin.Context) {
		// placeholder for real metrics
		c.JSON(http.StatusOK, gin.H{"requests": 1234})
	})

	// API group
	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			// list users
			users.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"users": []gin.H{
						{"id": 1, "name": "alice"},
						{"id": 2, "name": "bob"},
					},
				})
			})

			// create user
			type createUserReq struct {
				Name string `json:"name" binding:"required"`
			}
			users.POST("", func(c *gin.Context) {
				var req createUserReq
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				// fake created user
				c.JSON(http.StatusCreated, gin.H{"id": 3, "name": req.Name})
			})

			// get user by id
			users.GET("/:id", func(c *gin.Context) {
				id := c.Param("id")
				c.JSON(http.StatusOK, gin.H{"id": id, "name": "example"})
			})
		}
	}

	// Admin group with simple basic-auth middleware example
	admin := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"admin": "password", // DON'T use hardcoded credentials in production
	}))
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			user := c.MustGet(gin.AuthUserKey).(string)
			c.JSON(http.StatusOK, gin.H{"message": "welcome to admin dashboard", "user": user})
		})
	}

	// Serve static files (example)
	r.Static("/static", "./public")

	return r
}
