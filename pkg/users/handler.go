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

		_, err = s.CreateUser(&user)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			jsonMessage, _ := json.Marshal(err)
			w.Write(jsonMessage)
			return
		}
		
		jsonMessage, _ := json.Marshal("Successful registration")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonMessage)
	}
}


func loginUser(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		json.NewDecoder(r.Body).Decode(&user)

		hashedPassword, err := s.GetUserByEmail(&user) 
		
		if err != nil {
			jsonMessage, _ := json.Marshal("Invalid username or password")
			w.WriteHeader(http.StatusBadRequest)
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
	
		token, err := claims.SignedString([]byte("blabla"))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadGateway)
			jsonMessage, _ := json.Marshal("Login Failed. Please try again later.")
			w.Write(jsonMessage)
		}

		expiration := time.Now().Add(365 * 24 * time.Hour)
        cookie    :=    http.Cookie{Name: "jwtcookie",Value: token,Expires: expiration}
        http.SetCookie(w, &cookie)

		jsonMessage, _ := json.Marshal("Successful Login")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonMessage)
	}
} 
