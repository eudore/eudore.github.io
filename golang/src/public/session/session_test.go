package session;

import (
	"testing"
	"public/log"
	_ "public/cache/memcache"
	_ "public/session/store"
	"public/session"
)

func TestCache(t *testing.T) {
	sessionConfig := &session.ManagerConfig{CookieName: "token",EnableSetCookie: true, Gclifetime: 3600, Maxlifetime: 3600, Secure: true, CookieLifeTime: 3600, ProviderConfig: `{"store": "memcache","name": "token","conn":"127.0.0.1:12003"}`}
	globalSessions, err := session.NewManager("cache", sessionConfig)
	log.Info(err)
	go globalSessions.GC()
}