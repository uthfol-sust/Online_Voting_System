package routes

import (
	"fmt"
	"net/http"
	"pollvoting/pkg/controllers"
	"pollvoting/pkg/middleware"
)

func Router(router *http.ServeMux, userC controllers.UserController) {
	manager := middleware.Manager{}
	manager.Use(middleware.CorsMiddleware)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	UserRouter(router,userC, manager)
	PollsRouter(router, manager)
}
