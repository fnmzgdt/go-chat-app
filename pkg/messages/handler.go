package messages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func sendMessageToUser(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request)  {
		var message Message
		json.NewDecoder(r.Body).Decode(&message)
		message.Date = time.Now().Unix()

		if message.ThreadId == 0 {
			fmt.Println("no thread id in req body")
			threadId, err := s.GetUserThread(message.FromId, message.ToId)
			if err != nil {
				fmt.Println(err)
			}
			if threadId == 0 {
				fmt.Println("thread id doesnt exist between these users")
				var thread Thread
				thread.Type = 0

				res, err := s.CreateThread(&thread, message.Date) //passes

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
		//if there is a thread in the req body
		s.CreateMessage(&message)
		jsonMessage, _ := json.Marshal("Successfuly sent message")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonMessage)
		return
	}
}