
# Go Health - health-go

1. Run application

    ```bash
    go run main.go
    ```

1. Check application health

    - browser
        ```bash
        open http://localhost:8080/healthz
        ```

    - curl
        ```bash
        curl localhost:8080/healthz
        ```

    - curl + jq
        ```bash
        curl localhost:8080/healthz | jq
        ```

    - httpie
        ```bash
        http GET localhost:8080/healthz
        ```

## Links

- https://github.com/hellofresh/health-go

### available checks

- https://github.com/hellofresh/health-go/tree/master/checks
