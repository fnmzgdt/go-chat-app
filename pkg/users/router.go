package users

import "github.com/go-chi/chi"

func UsersRoutes(s Service) *chi.Mux {
	router := chi.NewRouter()
	router.Post("/login", loginUser(s))
	router.Post("/register", createUser(s))
	router.Post("/profile", changeProfilePicture(s))
	router.Delete("/profile", deleteProfilePicture(s))
	return router
}
