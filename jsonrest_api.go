package main

import (
	"encoding/json"
	"fmt"
	"github.com/function61/eventhorizon-exampleapp-go/events"
	"github.com/function61/eventhorizon-exampleapp-go/schema"
	"github.com/function61/eventhorizon-exampleapp-go/types"
	"github.com/function61/eventhorizon/pusher/pushlib/writerproxyclient"
	"github.com/function61/eventhorizon/util/cryptorandombytes"
	wtypes "github.com/function61/eventhorizon/writer/types"
	"log"
	"net/http"
	"os"
	"time"
)

func (a *App) setupJsonRestApi() {
	wpc := writerproxyclient.New()

	hostname, _ := os.Hostname()

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

	http.Handle("/orders", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var orders []schema.Order
		err := a.db.All(&orders)
		if err != nil {
			panic(err)
		}

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(orders); err != nil {
			panic(err)
		}
	}))

	// example of raising an event from inside our application and sending it to
	// Event Horizon. it will eventually reach our Pusher endpoint and get committed to
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

	http.Handle("/command/place_order", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var data types.OrderPlacement
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var user schema.User

		// make sure the user exists
		if err := a.db.One("ID", data.User, &user); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		orderId := cryptorandombytes.Hex(8)

		nowSerialized := time.Now().Format("2006-01-02 15:04:05")

		orderCreated := &events.OrderCreated{
			Id:   orderId,
			User: user.ID,
			Ts:   nowSerialized,
		}

		evs := []string{orderCreated.Serialize()}

		for _, item := range data.LineItems {
			lineItemAdded := &events.LineItemAdded{
				Order:   orderId,
				Product: item.Product,
				Amount:  item.Amount,
				Ts:      nowSerialized,
			}

			evs = append(evs, lineItemAdded.Serialize())
		}

		_, err := wpc.Append(&wtypes.AppendToStreamRequest{
			Stream: "/example",
			Lines:  evs,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// so Feeder can detect when the server changes
		w.Header().Set("X-Instance", hostname)

		w.Write([]byte(orderId))
	}))
}
