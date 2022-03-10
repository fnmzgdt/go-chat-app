package messages

import (
	"encoding/json"
	"net/http"
)

func createMessage(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func (w http.ResponseWriter, r *http.Request)  {
		var message Message
		json.NewDecoder(r.Body).Decode(&message)

		s.SendMessage(&message)
		w.Write([]byte("sent message"))
	}
}