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

// implements PushAdapter interface from pushlib
type App struct {
	pushListener *pushlib.Listener
	db           *storm.DB
}

func NewApp() *App {
	// subscription ID basically refers to a group of streams that we want to
	// receive updates for. therefore, all applications that want to follow
	// exactly the same streams can have the same subscription ID (like multiple
	// instances of same app in high-availability mode).
	// in multi-tenant architectures if you split your tenants between servers/clusters:
	// - cluster 1 could subscribe to /tenants/1, /tenants/3, /tenants/5, ...
	// - cluster 2 could subscribe to /tenants/2, /tenants/4, /tenants/6 ...
	// 
	// and moving /tenant/54321 between clusters is simply unsubscribing cluster 1'
	// subscription from /tenant/54321 and subscribing cluster 2's subscription to /tenant/54321
	subscriptionId := "foo"

	db, err := storm.Open("/tmp/app.db")
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	a := &App{
		db: db,
	}

	// init pushlib. we give reference to an object (this - our app) that
	// implements "PushAdapter" interface (all methods in this file prefixed
	// "Push"), which pushlib calls to process incoming event streams
	a.pushListener = pushlib.New(
		subscriptionId,
		a)

	return a
}

func (a *App) Run() {
	// start serving data over REST for companies and users (our data model)
	a.setupJsonRestApi()

	// start Pusher child process which will push stream updates to our HTTP
	// endpoint. the child process automatically exits if the parent (= us) exits,
	// so 1:1 relationship with Pusher <=> app endpoint is kind of enforced.
	// 
	// this design means that you cannot have multiple instances of your app
	// running per server unless a) your app instances use different ports
	// b) you use Docker so all the instances have their own network namespace.
	go pushlib.StartChildProcess("http://127.0.0.1:8080/_pyramid_push")

	// sets up HTTP path "/_pyramid_push" for receiving pushes
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

// We keep track of the offsets for each stream we follow. This example app
// follows just two:
// 
// - /_subscriptions/foo
// - /foostream
func (a *App) PushGetOffset(stream string, tx interface{}) (string, bool) {
	txReal := tx.(*transaction.Tx)

	// We don't have to verify stream names because those are based on the
	// subscription and we are only subscribed to streams that we know we want to follow.

	offset := ""
	if err := txReal.Db.WithTransaction(txReal.Tx).Get("cursors", stream, &offset); err != nil {
		if err == storm.ErrNotFound {
			// if we don't yet have an offset stored, that instructs Pusher to
			// start reading from the stream beginning.
			return "", false
		}

		// database read error?
		panic(err)
	}

	// ok we had an offset stored. pushlib asserts that new pushes continue from
	// this offset, guaranteeing us exactly-once delivery (no missed events, no re-processes)
	return offset, true
}

// called at end of stream processing to set the offset-in-stream from which we
// expect the next Push to start at.
func (a *App) PushSetOffset(stream string, offset string, tx interface{}) error {
	txReal := tx.(*transaction.Tx)

	if err := txReal.Db.WithTransaction(txReal.Tx).Set("cursors", stream, offset); err != nil {
		return err
	}

	return nil
}

// pushlib calls this to wrap all of the stream processing operations
// (GetOffset, HandleEvent, SetOffset) inside a transaction.
func (a *App) PushWrapTransaction(run func(interface{}) error) error {
	// use Bolt to get a write transaction (exclusively locked)
	err := a.db.Bolt.Update(func(tx *bolt.Tx) error {
		// wrap the DB (Storm perspective) and transaction (Bolt perspective)
		// in a struct that pushlib passes to our concrete handlers
		wrappedTransaction := &transaction.Tx{a.db, tx}

		// "run" is a pushlib callback that receives the tx and starts calling
		// the stream processing operations. if any of those return error, the
		// process is aborted and that error is returned from "run", and we return
		// the error from our callback (this code block) to Bolt, which in turn
		// rolls back the transaction on error or commits if no error
		return run(wrappedTransaction)
	})

	// ok transaction was either committed or rolled back, but this error is to let
	// pushlib know the final result to return over HTTP to the caller (probably Pusher)
	return err
}

func main() {
	receiver := NewApp()
	receiver.Run()

	cli.WaitForInterrupt()
}
