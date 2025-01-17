package main

import "log"
import "net/http"

func home(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("<h1>Hello from Exercise</h1>"))
}

func user(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("<h1>User list</h1>"))
}

func userDetail(w http.ResponseWriter, r *http.Request){
    w.Write([]byte("<h1>User detail</h1>"))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/{$}", home)
    mux.HandleFunc("/users", user)
    mux.HandleFunc("/users/detail", userDetail)
    log.Println("Starting server on :3000")
    err := http.ListenAndServe(":3000", mux)
    log.Fatal(err)
}