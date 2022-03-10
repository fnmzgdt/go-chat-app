package messages

import "github.com/go-chi/chi"

func MessagesRoutes(s Service) *chi.Mux { 
	router := chi.NewRouter()
	router.Post("/post", createMessage(s))
	return router
}