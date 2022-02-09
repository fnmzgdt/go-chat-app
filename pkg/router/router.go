package router

import (
	"fmt"
	"log"
	"net/http"
	"project/pkg/database"
	"project/pkg/users"

	"github.com/go-chi/chi"
)

func StartServer() *chi.Mux {
	fmt.Println("Starting server")

	r := database.SetupCassandraConnection()
	us := users.NewService(r)

	router := chi.NewRouter()
	router.Mount("/api/users", users.UsersRoutes(us))

	fmt.Println("Server is listening on PORT 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

	return router
}