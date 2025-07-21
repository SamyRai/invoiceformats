# Testing

Run all tests:

```sh
go test ./...
```

- Unit tests cover all major logic, including edge cases and error handling.
- See `providers/zugferd/builder_test.go`, `pkg/render/render_test.go`, etc.
- Test logger available in `testutils/`.
