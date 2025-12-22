package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/LittleJake/server-monitor-go/internal/controller"
	"github.com/LittleJake/server-monitor-go/internal/controller/api"
	"github.com/LittleJake/server-monitor-go/internal/middleware"
	"github.com/LittleJake/server-monitor-go/internal/util"

	"github.com/gin-gonic/gin"
)

// SetupRouter builds and returns a gin.Engine with example routes and middleware.
func SetupRouter() *gin.Engine {
	// parse bool from environment variable `WEB.IS_DEBUG`.
	// Accepts common boolean string values like: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
	isDebug := false
	if b := util.GetEnvBool("IS_DEBUG", false); b {
		isDebug = b
	}

	println("Web Debug Mode:", isDebug)
	if isDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.SetFuncMap(template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"trim":  strings.TrimSpace,
		// "json": func(v any) template.JS {
		// 	b, _ := json.Marshal(v)
		// 	return template.JS(b)
		// },
		"iconURL": func(v any) string {
			url := "https://cdnjs.cloudflare.com/ajax/libs/simple-icons/14.3.0/"
			icon := []string{
				"redhat",
				"centos",
				"ubuntu",
				"debian",
				"windows",
				"intel",
				"amd",
				"android",
				"qualcomm",
				"mediatek",
				"alpine linux",
				"arm",
				"openwrt",
				"qemu",
				"raspberrypi",
				//last
				"linux",
			}

			for _, icon := range icon {
				if strings.Contains(strings.ToLower(v.(string)), icon) {
					return fmt.Sprintf("%s%s.svg", url, strings.ReplaceAll(strings.ReplaceAll(icon, " ", ""), ".", "dot"))
				}
			}

			return "https://cdnjs.cloudflare.com/ajax/libs/simple-icons/14.3.0/linux.svg"
		},
		"iconName": func(v any) string {
			icon := []string{
				"redhat",
				"centos",
				"ubuntu",
				"debian",
				"windows",
				"intel",
				"amd",
				"android",
				"qualcomm",
				"mediatek",
				"alpine linux",
				"arm",
				"openwrt",
				"qemu",
				"raspberrypi",
				//last
				"linux",
			}

			for _, icon := range icon {
				if strings.Contains(strings.ToLower(v.(string)), icon) {
					return icon
				}
			}

			return "linux"
		},
		"iconColor": func(v any) string {
			icon := map[string]string{
				"redhat":       "#EE0000",
				"centos":       "#262577",
				"ubuntu":       "#E95420",
				"debian":       "#A81D33",
				"windows":      "#0078D6",
				"intel":        "#0071C5",
				"amd":          "#ED1C24",
				"android":      "#3DDC84",
				"qualcomm":     "#3253DC",
				"mediatek":     "#EC9430",
				"alpine linux": "#0D597F",
				"arm":          "#0091BD",
				"openwrt":      "#00B5E2",
				"qemu":         "#FF6600",
				"raspberrypi":  "#A22846",

				//last
				"linux": "#FCC624",
			}

			for name, color := range icon {
				if strings.Contains(strings.ToLower(v.(string)), name) {
					return color
				}
			}

			return "#FCC624"
		},
		"datetime": func(v any) string {
			if reflect.TypeOf(v) == reflect.TypeOf(int64(0)) {
				return time.Unix(v.(int64), 0).Format("2006-01-02 15:04:05")
			}

			if reflect.TypeOf(v) == reflect.TypeOf("") {
				i, _ := strconv.ParseFloat(v.(string), 64)
				return time.Unix(int64(i), 0).Format("2006-01-02 15:04:05")
			}

			return ""
		},

		"sizeFormat": func(v any) string {
			size, _ := strconv.ParseFloat(v.(string), 64)
			//input MB
			if size > 1024*1024*1024 {
				return fmt.Sprintf("%.2f PB", size*1.0/1024/1024/1024)
			}

			if size > 1024*1024 {
				return fmt.Sprintf("%.2f TB", size*1.0/1024/1024)
			}

			if size > 1024 {
				return fmt.Sprintf("%.2f GB", size*1.0/1024)
			}

			return fmt.Sprintf("%.2f MB", size)
		},
		"hash": func(v any) string {
			return fmt.Sprintf("%x", md5.Sum([]byte(v.(string))))
		},
	})

	// Built-in middleware
	r.Use(gin.Logger())
	r.Use(middleware.ServerDataMiddleware())
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
	r.NoRoute(controller.Error.NoRouteError)
	r.NoMethod(controller.Error.NoMethodError)

	r.Use(middleware.CORS)

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
		_api.GET("/thermal/:uuid", api.Thermal.Get)
		_api.GET("/report/:uuid", api.Report.Get)
		_api.GET("/battery/:uuid", api.Battery.Get)
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
		"admin": "password", // DON"T use hardcoded credentials in production
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
