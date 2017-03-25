package main

import (
	"encoding/json"
	"fmt"
	"github.com/function61/pyramid-exampleapp-go/events"
	"github.com/function61/pyramid-exampleapp-go/schema"
	"github.com/function61/pyramid/pusher/pushlib/writerproxyclient"
	wtypes "github.com/function61/pyramid/writer/types"
	"log"
	"net/http"
	"time"
)

func (a *App) setupJsonRestApi() {
	wpc := writerproxyclient.New()

	http.Handle("/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var users []schema.User
		err := a.db.All(&users)
		if err != nil {
			panic(err)
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(users); err != nil {
			panic(err)
		}
	}))

	http.Handle("/companies", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var companies []schema.Company
		err := a.db.All(&companies)
		if err != nil {
			panic(err)
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(companies); err != nil {
			panic(err)
		}
	}))

	// example of raising an event from inside our application and sending it to
	// Pyramid. it will eventually reach our Pusher endpoint and get committed to
	// our database.
	//
	// NOTE: normally this command would take in JSON but I was feeling lazy and
	//       I just used URL params...
	http.Handle("/command/change_user_name", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("user")
		var user schema.User

		// make sure the user exists
		if err := a.db.One("ID", userId, &user); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Printf("App: changing username for previous name = %s", user.Name)

		// create an event
		userNameChanged := &events.UserNameChanged{
			UserId:  user.ID,
			NewName: r.URL.Query().Get("new_name"),
			Reason:  "just an example",
			Ts:      time.Now().Format("2006-01-02 15:04:05"),
		}

		// append it to Pyramid
		output, err := wpc.Append(&wtypes.AppendToStreamRequest{
			Stream: "/example",
			Lines:  []string{userNameChanged.Serialize()},
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte(fmt.Sprintf("OK; offset = %s", output.Offset)))
	}))
}
