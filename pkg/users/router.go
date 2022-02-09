package users

import "github.com/go-chi/chi"

func UsersRoutes(s Service) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/get", getHandler)
	router.Post("/get", postHandler(s))
	return router
}