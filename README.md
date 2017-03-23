What is this
------------

Demonstration for running applications on top of [Pyramid](https://github.com/function61/pyramid).

This tutorial assumes that you're coming from/are familiar with the
[Quickstart guide](https://github.com/function61/pyramid/tree/master/docs/quickstart.md)


Architecture of this app
------------------------

- Golang, Linux, Docker
- Embedded database (BoltDB)
- ORM layer (Storm) on top of BoltDB. Normally I wouldn't recommend using an ORM
  but it was the least amount of code to persist data to make this demo work.


Setup streams etc.
------------------

Our app will listen on events in `/foostream`.
[Enter Pyramid CLI](https://github.com/function61/pyramid/blob/master/docs/enter-pyramid-cli.md)
and create the stream:

```
$ pyramid stream-create /foostream
```

We'll need to create a subscription (`foo`) and subscribe it to `/foostream`:

```
$ pyramid stream-create /_subscriptions/foo
$ pyramid stream-subscribe /foostream foo
```

Now when taking a peek at our subscription stream, we should see the tip of the
`/foostream` having been advertised:

```
$ pyramid stream-liveread /_subscriptions/foo:0:0 10
.Created {"subscription_ids":[],"ts":"2017-03-22T14:55:49.557Z"}
.SubscriptionActivity {"activity":["/foostream:0:135:1.2.3.4"],"ts":"2017-03-22T14:55:59.597Z"}
```

Now, let's enter some The Office -themed sample data into the stream so our
database will not be empty:

```
$ pyramid stream-appendfromfile /foostream example-dataimport/import.txt
2017/03/22 14:55:59 Appending 21 lines
2017/03/22 14:55:59 Done. Imported 21 lines in 135.305288ms.
```

If you'd now liveread the foo subscription again, we'd notice that there are
new notifications on the stream:

```
$ pyramid stream-liveread /_subscriptions/foo:0:0 10
.Created {"subscription_ids":[],"ts":"2017-03-22T14:55:49.557Z"}
.SubscriptionActivity {"activity":["/foostream:0:135:1.2.3.4"],"ts":"2017-03-22T14:55:59.597Z"}
.SubscriptionActivity {"activity":["/foostream:0:2417:1.2.3.4"],"ts":"2017-03-22T14:56:04.601Z"}
.SubscriptionActivity {"activity":["/foostream:0:2535:1.2.3.4"],"ts":"2017-03-22T15:16:00.62Z"}
```

Now you understand the mechanism for how Pusher knows which streams to push to
the endpoint - it just monitors the subscription stream which subscribed streams
have stuff the endpoint should be aware of! There's a bit more but that's the basic idea.

You can now exit from the CLI.


Now start the example application
---------------------------------

First, build the application:

```
$ docker build -t pyramid-exampleapp-go .
```

Then, run your application. Like with the Pyramid CLI, you have to specify the
`STORE` so the Pusher component can push events to your application:

```
$ docker run -it --rm -e STORE='...' pyramid-exampleapp-go
2017/03/22 15:04:37 App: listening at :8080
2017/03/22 15:04:37 pusherchild: starting
2017/03/22 15:04:37 configfactory: downloading discovery file
...
2017/03/22 15:04:37 Pusher: reached the top for /_subscriptions/foo
```

Ok it's succesfully started.


Interacting with the example application
----------------------------------------

This service has an internal projection of all the events it saw, i.e. it has
a database and you can query the database via its REST endpoint:

```
$ curl http://localhost:8080/users
[
...
   {
        "ID": "e1dd2e26",
        "Name": "Kelly Kapoor",
        "Company": "629cfead"
    }
]
```

We can modify the data by submitting new events. Normally the service itself
would probably have a UI with forms on how to change the data, which would post
the changes as events to Pyramid and Pyramid would notify your app of those events.

But here's the interesting bit: we can skip that application altogether, and use
Pyramid to directly change Kelly's (`id=e1dd2e26`) name!

If you look at [events/usernamechanged.go](events/usernamechanged.go), you'll
learn that we can do this (in the CLI):

```
$ pyramid stream-append /foostream 'UserNameChanged {"user_id": "e1dd2e26", "new_name": "Kelly Kaling", "reason": "Married", "ts": "2017-03-22 00:00:00"}'
```

And now inspect the data:

```
$ curl http://localhost:8080/users
[
...
   {
        "ID": "e1dd2e26",
        "Name": "Kelly Kaling",
        "Company": "629cfead"
    }
]
```


TODO
----

- demonstration of high availability setup
- demonstration of zero-downtime migration of a live service
