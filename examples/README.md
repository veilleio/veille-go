# Examples

```sh
export VEILLE_API_KEY="YOUR_API_KEY"
go run ./examples/json
```

- `json/` — call the featured endpoint and print the JSON body.
- `request/` — call any path with `Request()` and handle API errors.
- `wait-for-job/` — start an async job and wait for its result.
