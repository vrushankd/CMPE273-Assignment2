package main

import (
	// Controllers are kept in this package
	"controllers"
	// Standard library packages
	"net/http"
	// Third party packages
	"github.com/julienschmidt/httprouter"
)

func main() {
	// Instantiate a new router
	router := httprouter.New()

	// Get a controller instance
	controller := controllers.NewUserController()

	// Add handlers
	router.GET("/locations/:id", controller.GetLocation)
	router.POST("/locations", controller.CreateLocation)
	router.PUT("/locations/:id", controller.UpdateLocation)
	router.DELETE("/locations/:id", controller.DeleteLocation)

	// Expose the server at port 3000
	http.ListenAndServe("localhost:3000", router)
}
