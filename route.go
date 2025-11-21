package main

import (
	"net/http"
	"time"

	"github.com/LittleJake/server-monitor-go/controller"
	"github.com/LittleJake/server-monitor-go/controller/api"

	"github.com/gin-gonic/gin"
)

// SetupRouter builds and returns a gin.Engine with example routes and middleware.
func SetupRouter() *gin.Engine {
	r := gin.New()

	// Built-in middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

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

	r.GET("/", controller.Index.index)
	// Info routes
	r.GET("/info/:uuid", controller.Index.info)

	// API group
	_api := r.Group("/api")
	{
		_api.GET("/cpu/:uuid", api.Cpu.get)
		_api.GET("/memory/:uuid", api.Memory.get)
		_api.GET("/disk/:uuid", api.Disk.get)
		_api.GET("/network/:uuid", api.Network.get)
		_api.GET("/io/:uuid", api.IO.get)
		_api.GET("/ping/:uuid", api.Ping.get)
		_api.GET("/swap/:uuid", api.Swap.get)
		_api.GET("/thermal/:uuid", api.Thermal.get)
		_api.GET("/report/:uuid", api.Report.get)
	}

	// Public routes
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "uptime": time.Now().UTC()})
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
