package global

import (
	"strings"
	"net/http"
)


func AuthServeHTTP(w http.ResponseWriter, r *http.Request) {
	if(strings.HasPrefix(r.URL.Path,"/file/")) {

	}
}