package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

/*
func getHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Helloworld")
}
*/

func createUser(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		json.NewDecoder(r.Body).Decode(&user)
		
		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if err != nil {
			log.Fatal(err)
		}
	
		user.Password = string(password[:])

		s.CreateUser(&user)

		fmt.Println("added user to cassandra")
		json.NewEncoder(w).Encode("added user")
	}
}

func loginUser(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		json.NewDecoder(r.Body).Decode(&user)

		hashedPassword, err := s.GetUserByEmail(&user) 
		
		if err != nil {
			//if there is no such email in the db return 400 invalid username or password
			jsonMessage, _ := json.Marshal(err)
			w.WriteHeader(err.Status)
			w.Write(jsonMessage)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonMessage, _ := json.Marshal("Invalid username or password")
			w.Write(jsonMessage)
			return
		}

		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer: user.Email,
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		})
	
		token, loginerr := claims.SignedString([]byte("blabla"))
		if loginerr != nil {
			fmt.Println(loginerr)
			w.WriteHeader(http.StatusBadRequest)
			jsonMessage, _ := json.Marshal("Login Failed. Please try again later.")
			w.Write(jsonMessage)
		}

		jsonMessage, _ := json.Marshal(token)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonMessage)
	}
}

