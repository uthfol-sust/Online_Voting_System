package routes

import (
	"net/http"
	"pollvoting/pkg/controllers"
	"pollvoting/pkg/middleware"
)

func UserRouter(router *http.ServeMux, controllers controllers.UserController, manager middleware.Manager) {
	router.Handle("POST /users",
		manager.With(
			http.HandlerFunc(controllers.SingUp),
		))

	router.Handle("GET /users/{id}",
		manager.With(
			http.HandlerFunc(controllers.GetByID),
		))

	router.Handle("GET /users",
		manager.With(
			http.HandlerFunc(controllers.GetAll),
		))

	router.Handle("PUT /users/{id}",
		manager.With(
			http.HandlerFunc(controllers.Update),
		))

	router.Handle("DELETE /users/{id}",
		manager.With(
			http.HandlerFunc(controllers.Delete),
		))
}
