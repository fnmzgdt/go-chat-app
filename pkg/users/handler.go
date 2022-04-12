package users

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"project/pkg/errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func createUser(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user UserSignup
		json.NewDecoder(r.Body).Decode(&user)

		err := user.checkFields()
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}

		user.Password = string(password[:])

		_, err = s.CreateUser(&user)

		if err != nil {
			fmt.Println(err)
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		errors.RespondWithJSON(w, http.StatusOK, "message", "Successful registration")
	}
}

func loginUser(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user UserSignin
		json.NewDecoder(r.Body).Decode(&user)

		user.checkFields()

		hashedPassword, err := s.GetUserPasswordFromEmail(user.Email)
		if err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, "Invalid username or password")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password)); err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, "Invalid username or password")
			return
		}

		claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
			Issuer:    user.Email,
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		})

		token, err := claims.SignedString([]byte("blabla"))
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		expiration := time.Now().Add(365 * 24 * time.Hour)
		cookie := http.Cookie{Name: "jwtcookie", Value: token, Expires: expiration}
		http.SetCookie(w, &cookie)

		errors.RespondWithJSON(w, http.StatusOK, "message", "Successful login.")
	}
}

func changeProfilePicture(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(r.URL.Query().Get("userid"))

		if userId == 0 {
			errors.RespondWithError(w, http.StatusBadRequest, "userid param missing from request.")
			return
		}

		if r.Method != "POST" {
			errors.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, 3<<20)
		if err := r.ParseMultipartForm(3 << 20); err != nil {
			errors.RespondWithError(w, http.StatusUnprocessableEntity, "The uploaded file is too big. Please choose an file that's less than 3MB in size")
			return
		}

		file, _, err := r.FormFile("profile")
		if err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			errors.RespondWithError(w, http.StatusUnsupportedMediaType, "The provided file format is not allowed. Please upload a JPEG or PNG image")
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		defer file.Close()
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		ctx := r.Context()

		imageUrl, err := s.UploadProfileImage(ctx, userId, file)
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		_, err = s.UpdateProfileImage(ctx, userId, imageUrl)

		errors.RespondWithJSON(w, http.StatusOK, "message", "Successfully updated profile image.")
	}
}

func deleteProfilePicture(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(r.URL.Query().Get("userid"))

		if userId == 0 {
			errors.RespondWithError(w, http.StatusBadRequest, "userid param missing from request.")
			return
		}

		ctx := r.Context()

		err := s.DeleteProfileImage(ctx, userId)
		if err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		_, err = s.UpdateProfileImage(ctx, userId, "")
		if err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		errors.RespondWithJSON(w, http.StatusOK, "message", "Profile Image successfully deleted.")
	}
}