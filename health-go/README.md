
# Go Health - health-go

1. Run application

  ```bash
  go run main.go
  ```

1. Check application health

browser
  ```bash
  open http://localhost:8080/healtz
  ```

curl
  ```bash
  curl localhost:8080/healtz
  ```

curl + jq
  ```bash
  curl localhost:8080/healtz | jq
  ```

httpie
  ```bash
  http GET localhost:8080/healtz
  ```

## Links

- https://github.com/hellofresh/health-go
