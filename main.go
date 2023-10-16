package main

import "github.com/toggle-assignment/routes"

func main() {

	router := routes.SetupRoutes()

	router.Run("localhost:8080")
}
