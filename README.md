What is this
------------

Demonstration for running applications on top of Pyramid.

This tutorial assumes that you're coming from/are familiar with the
[Quickstart guide](https://github.com/function61/pyramid/tree/master/docs/quickstart.md)


Architecture of this app
------------------------

- Embedded database (BoltDB)
- ORM layer (Storm) on top of BoltDB. Normally I wouldn't recommend using an ORM
  but it was the least amount of code to persist data.


Setup streams etc.
------------------

Our app will listen on events in `/foostream`. Create it:

```
$ pyramid stream-create /foostream
```

We'll need to create a subscription (`foo`) and subscribe it to `/foostream`:

```
$ pyramid stream-create /_subscriptions/foo
$ pyramid stream-subscribe /foostream foo
```

Now, let's enter some sample data into the stream so our database will not be empty:

```
$ pyramid stream-appendfromfile /foostream example-dataimport/import.txt
2017/03/22 14:55:59 Appending 21 lines
2017/03/22 14:55:59 Done. Imported 21 lines in 135.305288ms.
```

Now start the service
---------------------

```
$ docker build -t pyramid-exampleapp-go .
$ docker run -it --rm -e STORE='...' pyramid-exampleapp-go
2017/03/22 15:04:37 App: listening at :8080
2017/03/22 15:04:37 pusherchild: starting
2017/03/22 15:04:37 configfactory: downloading discovery file
...
2017/03/22 15:04:37 Pusher: reached the top for /_subscriptions/foo
```

Ok it's succesfully started.


Interacting with the service
----------------------------

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
learn that we can do this:

```
$ pyramid stream-append /foostream 'UserNameChanged {"user_id": "e1dd2e26", "ts": "2017-03-22 00:00:00", "new_name": "Kelly Kaling", "reason": "Married"}'
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
