package global

import (
	"sync"
	"time"
	"strings"
	"strconv"
	"net/http"
	"golang.org/x/time/rate"
	"public/log"
)


func init() {
	go cleanupVisitors()
}

type gw struct {

}

func GetRealClientIP(r *http.Request ) string {
	xforward := r.Header.Get("X-Forwarded-For")
	if "" == xforward {
		return strings.SplitN(r.RemoteAddr, ":", 2)[0]
	}

	return strings.SplitN(string(xforward), ",", 2)[0]
}

func (g *gw) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	ip := GetRealClientIP(r)
	log.Info(ip," ",r.Method,": ",r.URL.Path)
	//echo(w,r)
	limiter := getVisitor(ip)
	if !limiter.Allow() {
		log.Info("限速：", ip)
		http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
		return
	}
	w.Header().Add("X-Request-Id","000")
	w.Header().Set("X-Runtime", strconv.FormatFloat(time.Since(t).Seconds(), 'f', -1, 64) )
	AuthServeHTTP(w,r)
	Router.ServeHTTP(w,r)
}


type Limitconifg struct{
	Rate	int
	Burst	int
}


type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}
// Change the the map to hold values of the type visitor.
var visitors = make(map[string]*visitor)
var mtx sync.Mutex

func getVisitor(ip string) *rate.Limiter {
	mtx.Lock()
	v, exists := visitors[ip]
	if !exists {
		mtx.Unlock()
		return addVisitor(ip)
	}
	// Update the last seen time for the visitor.
	v.lastSeen = time.Now()
	mtx.Unlock()
	return v.limiter
}
func addVisitor(ip string) *rate.Limiter {
	limiter := rate.NewLimiter(3, 10)
	mtx.Lock()
	// Include the current time when creating a new visitor.
	visitors[ip] = &visitor{limiter, time.Now()}
	mtx.Unlock()
	return limiter
}
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mtx.Lock()
		for ip, v := range visitors {
			if time.Now().Sub(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mtx.Unlock()
	}
}
