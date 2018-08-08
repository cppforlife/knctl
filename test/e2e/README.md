## End-to-end test

To run

```bash
# All tests
$ GOCACHE=off go test ./test/e2e/ -test.v

# Single file
$ GOCACHE=off go test ./test/e2e/e2e.go ./test/e2e/revisions_test.go -test.v

# Single test
$ GOCACHE=off go test ./test/e2e/ -test.v -run TestBuildFailed
```

See `./test/e2e/env.go` for required environment variables for some tests.
