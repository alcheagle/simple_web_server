package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func handler(w http.ResponseWriter, r *http.Request) {
    name, err := os.Hostname()
    if err != nil {
        panic(err)
    }
    fmt.Fprintf(w, "Hi there, I am %s", name)
}

func main() {
    http.HandleFunc("/", handler)
    log.Println("Starting server...")
    log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
