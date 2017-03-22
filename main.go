package main

import (
	"github.com/asdine/storm"
	"github.com/boltdb/bolt"
	"github.com/function61/pyramid/cli"
	"github.com/function61/pyramid-exampleapp-go/events"
	"github.com/function61/pyramid-exampleapp-go/transaction"
	"github.com/function61/pyramid/pusher/pushlib"
	"github.com/function61/pyramid/util/lineformatsimple"
	"log"
	"net/http"
)

// implements PushAdapter interface
type App struct {
	pushListener *pushlib.Listener
	db           *storm.DB
}

func NewApp() *App {
	subscriptionId := "foo"

	db, err := storm.Open("/tmp/listener.db")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	a := &App{
		db: db,
	}

	a.pushListener = pushlib.New(
		subscriptionId,
		a)

	return a
}

func (a *App) Run() {
	a.setupJsonRestApi()

	go pushlib.StartChildProcess("http://127.0.0.1:8080/_pyramid_push")

	// sets up HTTP App for receiving pushes
	a.pushListener.AttachPushHandler()

	srv := &http.Server{Addr: ":8080"}

	log.Printf("App: listening at :8080")

	if err := srv.ListenAndServe(); err != nil {
		// cannot panic, because this probably is an intentional close
		log.Printf("App: ListenAndServe() error: %s", err)
	}
}

// this is where all the magic happens. pushlib calls this function for every
// incoming event from Pyramid.
func (a *App) PushHandleEvent(eventSerialized string, tx interface{}) error {
	txReal := tx.(*transaction.Tx)

	// 'FooEvent {"bar": "input here"}'
	//     => eventType='FooEvent'
	//     => payload='{"bar": "input here"}'
	eventType, payload, err := lineformatsimple.Parse(eventSerialized)

	if err != nil {
		return err
	}

	if fn, fnExists := events.EventNameToApplyFn[eventType]; fnExists {
		return fn(txReal, payload)
	}

	log.Printf("App: unknown event: %s", eventSerialized)

	return nil
}

func (a *App) PushGetOffset(stream string, tx interface{}) (string, bool) {
	txReal := tx.(*transaction.Tx)

	offset := ""
	if err := txReal.Db.WithTransaction(txReal.Tx).Get("cursors", stream, &offset); err != nil {
		if err == storm.ErrNotFound {
			return "", false
		}

		// more serious error
		panic(err)
	}

	return offset, true
}

func (a *App) PushSetOffset(stream string, offset string, tx interface{}) error {
	txReal := tx.(*transaction.Tx)

	if err := txReal.Db.WithTransaction(txReal.Tx).Set("cursors", stream, offset); err != nil {
		return err
	}

	return nil
}

func (a *App) PushWrapTransaction(run func(interface{}) error) error {
	err := a.db.Bolt.Update(func(tx *bolt.Tx) error {
		return run(&transaction.Tx{a.db, tx})
	})

	return err
}

func main() {
	receiver := NewApp()
	receiver.Run()

	cli.WaitForInterrupt()
}
