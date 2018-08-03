package middlewares

import (
	"time"
	"strconv"
	"net/http"
	"public/log"
	"module/global"
	"module/auth/ram"
	"public/middlewares/rate"
	gzip "github.com/NYTimes/gziphandler"
)

type Middlewares struct {
	handles []http.Handler
}

func (g *Middlewares) AddHandle(h ...http.Handler) {
	g.handles = append(g.handles, h...)
}

func (g *Middlewares) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	log.Info(r.Method,": ",r.RequestURI )
	for _,i := range g.handles {
		i.ServeHTTP(w, r)
		if r.Method == "Deny" {
			break
		}
	}
	w.Header().Add("X-Request-Id","000")
	w.Header().Set("X-Runtime", strconv.FormatFloat(time.Since(t).Seconds(), 'f', -1, 64) )
}

var Md *Middlewares 

func init() {
	Md = &Middlewares{}
}

func Reload() error {
	Md.handles = make([]http.Handler, 0)
	ramhandle := &ram.Ram{
		DefaultPolicy:	[]string{"3", "4"},
	}
	Md.AddHandle(rate.NewRate(3, 10), ramhandle, gzip.GzipHandler(global.Router) )
	return nil
}
