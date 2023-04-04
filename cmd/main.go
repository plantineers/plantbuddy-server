package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/v1/hello", getHello)

	fmt.Println(http.ListenAndServe(":3333", nil))
}

func getHello(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	io.WriteString(w, "hello, world\n")
}
