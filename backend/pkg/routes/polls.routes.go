package routes

import (
	"net/http"
	"pollvoting/pkg/middleware"
	"pollvoting/pkg/controllers"
)

//API
//  POST /polls           -> create poll
// GET /polls            -> list polls
// GET /polls/{id}       -> get poll details
// DELETE /polls/{id}    -> delete poll (admin)
// PATCH /polls/{id}     -> update poll (admin)
// POST /polls/{id}/close -> manually close poll

func PollsRouter(router *http.ServeMux, controllers controllers.PollController, manager middleware.Manager) {

	router.Handle("POST /polls",
		manager.With(
			http.HandlerFunc(controllers.CreatePoll),
		))

	router.Handle("GET /polls",
		manager.With(
			http.HandlerFunc(controllers.GetPolls),
		))

	router.Handle("GET /polls/{id}",
		manager.With(
			http.HandlerFunc(controllers.GetPollDetails),
		))

	router.Handle("DELETE /polls/{id}",
		manager.With(
			http.HandlerFunc(controllers.DeletePoll),
		))

	router.Handle("PUT /polls/{id}",
		manager.With(
			http.HandlerFunc(controllers.UpdatePoll),
		))
}
