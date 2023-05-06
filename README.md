### Used

```shell
$ git clone https://github.com/chenpkg/balancer.git
$ cd ./balancer
$ go build
```

**run http test server**

```shell
$ cd ./httptest
$ go run server.go 8080 8081 8082 8083
```

**after run balancer**

```shell
$ cd ..
$ ./balancer
Schema: http
Port: 8000
Health Check: true
Location:
        Route: /
        Proxy Pass: [http://127.0.0.1:8080 http://127.0.0.1:8081 http://127.0.0.1:8082 http://127.0.0.1:8083 http://127.0.0.1:8084]
        Mode: round-robin

        Route: /random
        Proxy Pass: [http://127.0.0.1:8082 http://127.0.0.1:8083 http://127.0.0.1:8084]
        Mode: random
```

**and run http client**

```shell
$ cd ./httptest
$ go run client.go
127.0.0.1:8080
127.0.0.1:8081
```