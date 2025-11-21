package main

// optional: minimal main to run the router
func main() {
	r := SetupRouter()
	// listen and serve on 0.0.0.0:8080
	print("starting server.")
	_ = r.Run(":8080")
}
