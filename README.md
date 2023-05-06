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
```

**and run http client**

```shell
$ cd ./httptest
$ go run client.go
127.0.0.1:8080
```