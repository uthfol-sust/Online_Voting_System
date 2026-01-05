package routes

import (
	"net/http"
	"pollvoting/pkg/controllers"
	"pollvoting/pkg/middleware"
)

func VoteRoutes(router *http.ServeMux, controller controllers.CastVoteController, manager middleware.Manager) {
	router.Handle("POST /votes",
		manager.With(
			http.HandlerFunc(controller.CastVote),
			middleware.AuthMiddleware,
		))

	router.Handle("GET /votes/result/{pollID}",
		manager.With(
			http.HandlerFunc(controller.PollResults),
			middleware.AuthMiddleware,
		))
}

