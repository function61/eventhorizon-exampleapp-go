What is this
------------

Demonstration for running applications on top of
[Event Horizon](https://github.com/function61/eventhorizon).

This project assumes that you're coming from/are familiar with the
[Quickstart guide](https://github.com/function61/eventhorizon/tree/master/docs/quickstart.md).

The meat of the project is in [main.go](main.go) which is well documented and
you should read it.


Architecture of this app
------------------------

- Golang, Linux, Docker
- Embedded database (BoltDB)
- ORM layer (Storm) on top of BoltDB. Normally I wouldn't recommend using an ORM
  but it was the least amount of code to persist data to make this demo work.


Tutorials
---------

It is recommended that you read these in order.

1. [Setting up & introduction](docs/setting-up-and-introduction.md)
2. [Live migration tutorial](docs/live-migration-tutorial.md)


TODO
----

While under heavy writing:

- demonstration of high availability setup
- demonstration of zero-downtime migration of a live service to 100 % different tech stack
- demonstration of pause + database structure migration + resume
