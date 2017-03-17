package main

import (
	"net/http"
)

func test(w http.ResponseWriter, r *http.Request){
	w.Header().Add("Content-Type", "text/html")
	w.Write([]byte("Testing  Page"))

}
