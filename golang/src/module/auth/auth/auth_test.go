package auth

import (
	"time"
	"testing"
	"module/auth/auth"
)


var ts string

func TestAuth(t *testing.T) {
	to,err := auth.Auth("/file/",auth.GET | auth.POST ,time.Now().Add(100 * time.Second).Unix())
	t.Log(to,err,len(to))
	ts = to
}

func TestValid(t *testing.T) {
	path := "/file/weer/public/ss"
	t.Log(path,auth.Valid(ts,path,auth.GET))
	t.Log(path,auth.Valid(ts,path,auth.DELETE))
}