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

	router.Handle("POST /login",
		manager.With(
			http.HandlerFunc(controllers.Login),
		))

	router.Handle("POST /refresh",
		manager.With(
			http.HandlerFunc(controllers.RefreshToken),
		))

	router.Handle("GET /users/{id}",
		manager.With(
			http.HandlerFunc(controllers.GetByID),
			middleware.AuthMiddleware,
		))

	router.Handle("GET /users",
		manager.With(
			http.HandlerFunc(controllers.GetAll),
			middleware.AuthMiddleware,
		))

	router.Handle("PUT /users/{id}",
		manager.With(
			http.HandlerFunc(controllers.Update),
			middleware.AuthMiddleware,
		))

	router.Handle("DELETE /users/{id}",
		manager.With(
			http.HandlerFunc(controllers.Delete),
			middleware.AuthMiddleware,
		))
}
