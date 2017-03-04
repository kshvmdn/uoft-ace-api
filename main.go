package main

import (
	"fmt"
	"os"
)

func main() {
	a := App{}

	a.Initialize(
		os.Getenv("REDIS_PORT"),
		os.Getenv("REDIS_PASSWORD"),
		os.Getenv("REDIS_DB"))

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	a.Run(fmt.Sprintf(":%s", port))
}
