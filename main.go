package main

import (
	"net/http"
)

var collection = NewCollection()

func main() {
	go collection.run()
	http.HandleFunc("/ws", WSHandler)
	http.HandleFunc("/push", ApiHandler)
	http.ListenAndServe(":8001", nil)
}
