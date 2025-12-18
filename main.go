package main

import (
	"fmt"

	"github.com/LittleJake/server-monitor-go/internal/util"
)

// optional: minimal main to run the router
func main() {
	//load .env
	_ = util.LoadEnv()

	r := SetupRouter()
	// listen and serve on 0.0.0.0:8080
	r.LoadHTMLGlob("view/**/*.html")

	redisClient = NewRedisClient(
		fmt.Sprintf("%s:%s",
			util.GetEnv("REDIS_HOST", "127.0.0.1"),
			util.GetEnv("REDIS_PORT", "6379"),
		),
		util.GetEnv("REDIS_PASSWORD", ""),
		util.GetEnvInt("REDIS_DB", 0),
	)

	print("starting server.")
	_ = r.Run(":8888")
}
