package ram

import (
	"testing"
	"module/auth/ram"
)

func TestMatch(t *testing.T) {
	t.Log(ram.MatchStar("*","/sss"))
	t.Log(ram.MatchStar("/note/*","/note/"))
	t.Log(!ram.MatchStar("/note/*","/api/note/golang"))
	t.Log(ram.MatchStar("/note/*lang","/note/golang"))
	t.Log(!ram.MatchStar("/note/*langs","/note/golang/gets"))
	t.Log(ram.MatchStar("*.html","/note/golang.html"))
	t.Log(!ram.MatchStar("*.html","/note/golang.js"))
}