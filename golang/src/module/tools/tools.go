package tools;

import (
    "net/http"
	"module/global"
)

func Reload() error {
	global.Router.GetFunc("/tools/", tools)
	return nil
}


func tools(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("tools" + r.URL.Path))
}
