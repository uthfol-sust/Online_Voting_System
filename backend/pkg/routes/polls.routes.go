package routes

import (
	"net/http"
	"pollvoting/pkg/middleware"
)

func PollsRouter(router *http.ServeMux,manager middleware.Manager){
	// router.Handle("POST /polls", manager.With(
	// 	http.HandlerFunc(),
	// 	middleware.AuthMiddleware,
	// ))
}