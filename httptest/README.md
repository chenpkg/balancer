### 测试程序

### server.go

```shell
$ go run server.go 8080 8081 8082 8083
```

启动4个 http 服务端，监听 8080-8083 端口

再另一个命令行窗口再次启动 8084 端口，不断启动和停用可观察到负载均衡器的健康检测机制。

```shell
$ go run server.go 8084
```

### client.go

```shell
$ go run client.go
```

启动 http 客户端连接到负载均衡器端口，命令行将显示访问到的端口。客户端并不会立即退出，每按一次回车将发起一次 http 请求，将观察到访问的端口变化。