package messages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func sendMessageToThread(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request)  {
		var message MessagePost
		json.NewDecoder(r.Body).Decode(&message)
		if err := message.validate(); err != nil {
			respondWithError(w, http.StatusBadRequest, "Message can't be empty")
			return
		}
		message.Date = time.Now().Unix()

		/*
		if message.Thread.Id == 0 {
			
			fmt.Println("no thread id in request body")
			threadId, err := s.GetUserThread(message.FromId, message.ToId)
			if err != nil {
				fmt.Println(err)
			}
			
			if threadId == 0 {
				fmt.Println("thread id doesnt exist between these users")
				var thread Thread
				thread.CreatorId = 0
				thread.Type = 0
				thread.CreatedAt = message.Date
				res, err := s.CreateThread(&thread) //passes

				LastInsertId,_ := res.LastInsertId()
				message.ThreadId = LastInsertId //passes

				//create thread for the user that sends first message
				_, err = s.CreateUserThread(message.FromId, message.ThreadId, message.Date, 1)
				if err != nil {
					fmt.Println(err)
				}

				//create thread for the user that recieves the first message
				_, err = s.CreateUserThread(message.ToId, message.ThreadId, message.Date, 0)
				if err != nil {
					fmt.Println(err)
				}

				_, err = s.CreateMessage(&message)
				if err != nil {
					fmt.Println(err)
				}
				
				jsonMessage, _ := json.Marshal("Successfuly sent message")
				w.WriteHeader(http.StatusOK)
				w.Write(jsonMessage)
				return
			} 
			
			// if the thread exist assign it to the struct
			message.ThreadId = threadId
			s.CreateMessage(&message)
			jsonMessage, _ := json.Marshal("Successfuly sent message")
			w.WriteHeader(http.StatusOK)
			w.Write(jsonMessage)
			return
		}
		*/
		//if there is a thread in the req body
		err,_ := s.CreateMessage(&message)
		if err != nil {
			fmt.Println(err)
			respondWithError(w, http.StatusBadGateway, "DB error")
			return
		}

		//update threads
		for _, id := range message.Thread.Users{
			var seen uint8
			seen = 0
			if id == message.FromId {
				seen = 1
			}
			_, err := s.UpdateUserThread(id, message.Thread.Id, seen)
			if err != nil {
				fmt.Println(err)
				respondWithError(w, http.StatusBadGateway, "DB error")
				return
			}
		}

		jsonMessage, _ := json.Marshal("Successfuly sent message")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonMessage)
		return
	}
}

func createGroupThread(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var thread Thread
		json.NewDecoder(r.Body).Decode(&thread)
		thread.Type = 1
		thread.CreatedAt = time.Now().Unix()

		res, err := s.CreateThread(&thread)
		if err != nil {
			respondWithError(w, http.StatusBadGateway, "DB error")
			return
		}
		threadid,_ := res.LastInsertId()

		for _, id := range thread.Users{
			var seen uint8
			seen = 0
			if id == thread.CreatorId {
				seen = 1
			}
		
			_, err = s.CreateUserThread(id, threadid, thread.CreatedAt, seen)
			if err != nil {
				fmt.Println(err)
				respondWithError(w, http.StatusBadGateway, "DB error")
				return
			}
		}

		jsonMessage, _ := json.Marshal("Successfuly created groupchat")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonMessage)
		return
	}
} 

func getMessagesFromThread(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		threadid,_ := strconv.Atoi(r.URL.Query().Get("threadid"))
		userid,_ := strconv.Atoi(r.URL.Query().Get("userid"))
		//get latest messages in thread
		res, err := s.GetMessagesFromThread(threadid)
		if err != nil {
			fmt.Println(err)
			respondWithError(w, http.StatusBadGateway, "DB error")
			return
		}

		_, err = s.UpdateUserThread(userid, threadid, 1)
		if err != nil {
			fmt.Println(err)
			respondWithError(w, http.StatusBadGateway, "DB error")
			return
		}

		response, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	}
}

func getLatestThreads(s Service) (func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		userid, _ := strconv.Atoi(r.URL.Query().Get("userid"))

		res, err := s.getLatestThreads(userid)
		if err != nil {
			fmt.Println(err)
			respondWithError(w, http.StatusBadGateway, "DB error")
			return
		}
		
		response, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	}
}