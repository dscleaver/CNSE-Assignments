## Voter API

This is an implementation of the voter API. All tests depend on the server running.

```
make build up
```

This defaults to starting the server on 0.0.0.0:1080 with a supporting redis cache. 
With the server up and running the tests can be run with:

```
make test
```

To stop the running api and cache simply run:

```
make stop
```

NOTE: The following additional extra credit items are implemented in this project:
* Added multi-platform build with `make build-multiplatform`