package tools;

import (
    "net/http"
	"public/router"
)

func init() {
	mux := router.Instance()
	mux.GetFunc("/tools/", tools)
}


func tools(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("tools" + r.URL.Path))
}
