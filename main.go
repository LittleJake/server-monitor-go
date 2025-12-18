package main

import "github.com/LittleJake/server-monitor-go/internal/util"

// optional: minimal main to run the router
func main() {
	//load .env
	_ = util.LoadEnv()

	r := SetupRouter()
	// listen and serve on 0.0.0.0:8080
	r.LoadHTMLGlob("view/**/*.html")
	print("starting server.")
	_ = r.Run(":8888")
}
