package tools;

import (
    "net/http"
	"public/router"
)

func init() {
	mux := router.Instance()
	mux.GetFunc("/tools/proxy/subscribe", subscribe)
	mux.GetFunc("/tools/proxy", proxy)
	mux.GetFunc("/tools/", tools)
}


func tools(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("tools"))
}
