package messages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/pkg/errors"
	"strconv"
	"time"
)

func sendMessageToThread(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var message MessagePost
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}

		senderId, _ := strconv.Atoi(r.URL.Query().Get("userid"))
		threadId, _ := strconv.Atoi(r.URL.Query().Get("threadid"))

		if senderId == 0 {
			errors.RespondWithError(w, http.StatusBadRequest, "userid param missing from request.")
			return
		}

		if threadId == 0 {
			errors.RespondWithError(w, http.StatusBadRequest, "threadid param missing from request.")
			return
		}

		usersInThread, err := s.GetUsersInThread(threadId)
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		message.FromId = senderId
		message.ThreadId = threadId
		message.Date = int(time.Now().Unix())
		message.ToUserIds = usersInThread

		if err := message.validate(); err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		_, err = s.CreateMessage(&message)
		if err != nil {
			fmt.Println(err)
			errors.RespondWithError(w, http.StatusInternalServerError, "DB error")
			return
		}

		_, err = s.UpdateUserThread(message.FromId, message.ThreadId, message.ToUserIds) //sets seen = 0 for all other users apart from the sender
		if err != nil {
			fmt.Println(err)
			errors.RespondWithError(w, http.StatusInternalServerError, "DB error")
			return
		}

		errors.RespondWithJSON(w, http.StatusOK, "message", "Successfuly sent message")
		return
	}
}

func createGroupThread(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var thread ThreadPost
		json.NewDecoder(r.Body).Decode(&thread)

		creatorId, _ := strconv.Atoi(r.URL.Query().Get("userid"))

		if creatorId == 0 {
			errors.RespondWithError(w, http.StatusBadRequest, "userid param missing from request.")
			return
		}

		thread.Type = 1
		thread.CreatorId = creatorId
		thread.CreatedAt = int(time.Now().Unix())

		err := thread.checkFields()
		if err != nil {
			errors.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := s.CreateThread(&thread)
		if err != nil {
			errors.RespondWithError(w, http.StatusInternalServerError, "DB error")
			return
		}
		threadid, _ := res.LastInsertId()

		_, err = s.UpdateUserThread(thread.CreatorId, int(threadid), thread.Users)
		if err != nil {
			fmt.Println(err)
			errors.RespondWithError(w, http.StatusInternalServerError, "DB error")
			return
		}

		errors.RespondWithJSON(w, http.StatusOK, "message", "Groupchat successfully created.")
		return
	}
}

func getMessagesFromThread(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		threadid, _ := strconv.Atoi(r.URL.Query().Get("threadid"))
		userid, _ := strconv.Atoi(r.URL.Query().Get("userid"))

		if threadid == 0 {
			errors.RespondWithError(w, http.StatusBadRequest, "threadid param missing from request.")
			return
		}

		if userid == 0 {
			errors.RespondWithError(w, http.StatusBadRequest, "userid param missing from request.")
			return
		}

		res, err := s.GetMessagesFromThread(threadid)
		if err != nil {
			fmt.Println(err)
			errors.RespondWithError(w, http.StatusInternalServerError, "DB error")
			return
		}

		_, err = s.UpdateUserThread(userid, threadid, []int{userid})
		//check rows matched ?? is it possible
		if err != nil {
			fmt.Println(err)
			errors.RespondWithError(w, http.StatusInternalServerError, "DB error")
			return
		}

		errors.RespondWithJSON(w, http.StatusOK, "payload", res)
		return
	}
}

func getLatestThreads(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userid, _ := strconv.Atoi(r.URL.Query().Get("userid"))

		if userid == 0 {
			errors.RespondWithError(w, http.StatusBadRequest, "userid param missing from request.")
			return
		}

		res, err := s.getLatestThreads(userid)
		if err != nil {
			fmt.Println(err)
			errors.RespondWithError(w, http.StatusInternalServerError, "DB error")
			return
		}

		errors.RespondWithJSON(w, http.StatusOK, "payload", res)
		return
	}
}
