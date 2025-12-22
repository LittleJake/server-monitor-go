package main

import (
	"log"

	"github.com/LittleJake/server-monitor-go/internal/util"
)

// optional: minimal main to run the router
func main() {
	// load .env
	_ = util.LoadEnv()

	// Setup Redis client and test connection
	if err := util.SetupRedis(); err != nil {
		log.Fatalf("Failed to setup Redis: %v", err)
	}
	defer func() {
		if err := util.CloseRedisClient(); err != nil {
			log.Printf("Error closing Redis: %v", err)
		}
	}()

	util.SetupCollectionCache()

	r := SetupRouter()
	// listen and serve on 0.0.0.0:8080
	r.LoadHTMLGlob("view/**/*.html")

	print("starting server on :8888")
	_ = r.Run("127.0.0.1:8888")
}
