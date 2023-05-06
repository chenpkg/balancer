package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func req(url string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}

func main() {
	url := "http://127.0.0.1:8000"
	req(url)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		req(url)
	}
}
