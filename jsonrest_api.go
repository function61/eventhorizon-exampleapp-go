package main

import (
	"encoding/json"
	"github.com/function61/pyramid-exampleapp-go/schema"
	"net/http"
)

func (a *App) setupJsonRestApi() {
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
}
