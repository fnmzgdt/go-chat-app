package users

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (user User)

func getHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Helloworld")
}

func postHandler(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&user)
		s.AddUser(&user)

		fmt.Println("added user to cassandra")
		json.NewEncoder(w).Encode("added user")
	}
}

