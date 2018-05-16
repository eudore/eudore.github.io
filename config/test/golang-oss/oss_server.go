package main;

import (
	"net/http"
)


func main() {
	http.HandleFunc("/policy", oss_policy)
	http.HandleFunc("/callback", oss_callback)
	http.ListenAndServe(":9090", nil)
}