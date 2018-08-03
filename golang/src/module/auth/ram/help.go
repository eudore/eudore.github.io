package ram

import (
	"strings"
	"net/http"
)

func MatchStar(patten, obj string) bool {
	ps := strings.Split(patten,"*")
	if len(ps) == 0 {
		return patten == obj
	}
	if !strings.HasPrefix(obj, ps[0]) {
		return false
	}
	for _,i := range ps {
		if i == "" {
			continue
		}
		pos := strings.Index(obj, i)
		if pos == -1 {
			return false
		}
		obj = obj[pos + len(i):]
	}
	return true
}


func GetRealClientIP(r *http.Request ) string {
	xforward := r.Header.Get("X-Forwarded-For")
	if "" == xforward {
		return strings.SplitN(r.RemoteAddr, ":", 2)[0]
	}

	return strings.SplitN(string(xforward), ",", 2)[0]
}