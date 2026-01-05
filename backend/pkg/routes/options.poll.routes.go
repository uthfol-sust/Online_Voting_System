package routes

import (
	"net/http"
	"pollvoting/pkg/controllers"
	"pollvoting/pkg/middleware"
)

func OptionsRoutes(router *http.ServeMux, controllers controllers.OptionController, manager middleware.Manager) {
	router.Handle("POST /options",
		manager.With(
			http.HandlerFunc(controllers.CreateOption),
			middleware.AuthMiddleware,
		))

	router.Handle("DELETE /options/{id}",
		manager.With(
			http.HandlerFunc(controllers.DeleteOption),
			middleware.AuthMiddleware,
		))
}
