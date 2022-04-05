package router

import (
	"fmt"
	"log"
	"net/http"
	"project/pkg/database"
	"project/pkg/messages"
	"project/pkg/users"
	"project/pkg/websocket"

	"github.com/go-chi/chi"
)

func StartServer() *chi.Mux {
	fmt.Println("Starting server")
	
	mysql, err := database.SetupMySQLConnection()

	if err != nil {
		fmt.Println(err)
	}

	us := users.NewService(mysql)
	ms := messages.NewService(mysql)

	hub := websocket.InitializeNewHub()
	go hub.SetupEventRouter()

	router := chi.NewRouter()
	router.HandleFunc("/chat", hub.ServeWs)
	router.Mount("/api/users", users.UsersRoutes(us))
	router.Mount("/api/messages", messages.MessagesRoutes(ms))

	fmt.Println("Server is listening on PORT 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

	return router
}