## Voter API

This is an implementation of the voter API. All tests depend on the server running.

```
go run main
```

This defaults to starting the server on 0.0.0.0:1080. With the server up and running
the tests can be run with:

```
go test ./...
```

NOTE: The following additional extra credit items are implemented in this project:
* Added json tags to returned structures
* Added put and delete for voters and polls
* Added meaningful data to the health endpoint