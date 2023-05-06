package main

import (
	"fmt"
	"net/http"
	"os"
)

type Handle struct {
	port string
}

func (h *Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "127.0.0.1:"+h.port)
}

func getAddress(port string) string {
	return "127.0.0.1:" + port
}

func httpServer(port string) {
	fmt.Printf("listener http port: %s\n", port)

	server := &http.Server{
		Addr:    getAddress(port),
		Handler: &Handle{port},
	}

	server.ListenAndServe()
}

func main() {
	ports := os.Args[1:]

	if len(ports) == 0 {
		ports = []string{"8080"}
	}

	for _, port := range ports {
		go httpServer(port)
	}

	select {}
}
