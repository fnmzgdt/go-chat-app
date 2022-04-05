package messages

import "github.com/go-chi/chi"

func MessagesRoutes(s Service) *chi.Mux { 
	router := chi.NewRouter()
	router.Post("/post", sendMessageToThread(s))
	router.Post("/thread", createGroupThread(s))
	router.Get("/thread", getMessagesFromThread(s))
	router.Get("/latestthreads", getLatestThreads(s))
	return router
} 