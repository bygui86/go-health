
# go-healthcheck
Example project to understand healthcheck go library

## Run
`INFO`: to better test this library, it's a good idea to avoid starting the DB
```shell
go run main.go
```

## Test
```shell
curl localhost:8080/healthcheck | jq
# or
http GET localhost:8080/healthcheck
```

---

## TODO
- [ ] test with a REST server
- [ ] test with a gRPC server
- [ ] test with a gRPC client

---

## Links

- https://github.com/etherlabsio/healthcheck
