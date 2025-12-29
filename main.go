package main

import (
	"fmt"
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
	util.SetupMapStringCache()
	util.SetupCollectionStatusCache()
	util.SetupDiskCache()

	r := SetupRouter()

	go util.CronJob()

	fmt.Printf("starting server on %s:%d", util.GetEnv("LISTEN_ADDRESS", "127.0.0.1"), util.GetEnvInt("LISTEN_PORT", 8888))
	_ = r.Run(fmt.Sprintf("%s:%d", util.GetEnv("LISTEN_ADDRESS", "127.0.0.1"), util.GetEnvInt("LISTEN_PORT", 8888)))
}
