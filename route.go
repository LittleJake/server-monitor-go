package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/LittleJake/server-monitor-go/internal/assets"
	"github.com/LittleJake/server-monitor-go/internal/controller"
	"github.com/LittleJake/server-monitor-go/internal/controller/api"
	"github.com/LittleJake/server-monitor-go/internal/middleware"
	"github.com/LittleJake/server-monitor-go/internal/util"
	"github.com/elliotchance/orderedmap/v3"
	"github.com/gin-contrib/i18n"

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

	// Built-in middleware
	r.Use(gin.Logger())
	r.Use(middleware.ServerDataMiddleware())
	r.Use(middleware.GinI18nLocalize())

	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"trim":  strings.TrimSpace,
		"iconURL": func(v any) string {
			url := "https://cdnjs.cloudflare.com/ajax/libs/simple-icons/14.3.0/"
			iconName := util.IconNameFormat(v)
			return fmt.Sprintf("%s%s.svg", url, strings.ReplaceAll(strings.ReplaceAll(iconName, " ", ""), ".", "dot"))
		},
		"iconName": util.IconNameFormat,
		"iconColor": func(v any) string {
			iconColor := map[string]string{
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

				// //last
				"linux": "#FCC624",
			}

			return iconColor[util.IconNameFormat(v)]
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

			units := []string{"MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

			i := 0

			for i = 0; size >= 1024 && i < len(units)-1; i++ {
				size /= 1024
			}

			return fmt.Sprintf("%.2f %s", size, units[i])
		},
		"hash": func(v any) string {
			return fmt.Sprintf("%x", md5.Sum([]byte(v.(string))))
		},
		"locale": i18n.GetMessage,
		"count": func(v any) int {
			if reflect.TypeOf(v) == reflect.TypeOf([]string{}) {
				return len(v.([]string))
			}
			if reflect.TypeOf(v) == reflect.TypeOf(map[string]interface{}{}) {
				return len(v.(map[string]interface{}))
			}
			if reflect.TypeOf(v) == reflect.TypeOf(orderedmap.OrderedMap[any, any]{}) {
				return v.(*orderedmap.OrderedMap[any, any]).Len()
			}
			return 0
		},
		"default": func(v any, d any) any {
			if v == nil || fmt.Sprintf("%s", v) == "" {
				return fmt.Sprintf("%s", d)
			}
			return fmt.Sprintf("%s", v)
		},
	}

	r.SetHTMLTemplate(template.Must(template.New("").Funcs(funcMap).ParseFS(assets.TemplatesFS, "templates/**/*.html")))

	// Use the default Recovery in debug mode so developers see panics.
	// In non-debug (production) mode, use a custom recovery that renders
	// a friendly error page instead of exposing stack traces.

	r.Use(gin.CustomRecovery(controller.Error.InternalServerError))

	// In non-debug mode, route unknown paths and methods to the error page.
	r.NoRoute(controller.Error.NoRouteError)
	r.NoMethod(controller.Error.NoMethodError)

	r.Use(middleware.CORS)

	r.GET("/", controller.Index.Index)

	// Info routes
	r.GET("/info/:uuid", controller.Index.Info)
	r.GET("/list/", controller.Index.List)

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
		_api.GET("/battery/:uuid", api.Battery.Get)

		_api.POST("/report/:uuid", api.Report.Set)
	}

	// Admin group with simple basic-auth middleware example
	// admin := r.Group("/admin", gin.BasicAuth(gin.Accounts{
	// 	"admin": "password", // DON"T use hardcoded credentials in production
	// }))
	// {
	// 	admin.GET("/dashboard", func(c *gin.Context) {
	// 		user := c.MustGet(gin.AuthUserKey).(string)
	// 		c.JSON(http.StatusOK, gin.H{"message": "welcome to admin dashboard", "user": user})
	// 	})
	// }
	static, _ := fs.Sub(assets.StaticFS, "static")
	// Serve static files (example)
	r.StaticFS("/static", http.FS(static))

	// Serve embedded favicon.ico
	r.GET("/favicon.ico", assets.ServeFavicon)

	r.GET("/manifest.json", assets.ServeManifest)

	return r
}
