package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

type Manager struct {
	gobalMiddlewares []Middleware
}

func (mngr *Manager) Use(middlewares ...Middleware) {
	mngr.gobalMiddlewares = append(mngr.gobalMiddlewares, middlewares...)
}

func (mngr *Manager) With(next http.Handler, middlewares ...Middleware) http.Handler {
	n := next

	for _, middleware := range middlewares {
		n = middleware(n)
	}

	for _, gobal := range mngr.gobalMiddlewares {
		n = gobal(n)
	}

	return n
}
