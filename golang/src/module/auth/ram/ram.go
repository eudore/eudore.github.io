package ram

import (
	"fmt"
	"strings"
	"net/http"
	"database/sql"
	"public/log"
	"public/token"
	"module/global"
)

var (
	stmtQueryRamPolicy *sql.Stmt
)

func Reload() error {
	stmtQueryRamPolicy = global.Stmt("SELECT `Type`,`Config`,`Condition`,`Effect` FROM tb_auth_ram_policy WHERE `PID`=?;")
	// log.Info(Match([]string{"1","0"},map[string]string{"url":"/file/golang", "method": "Put"}))
	// log.Info(Match([]string{"1"},map[string]string{"url":"/note/golang", "method": "Get"}))
	// log.Info(!Match([]string{"1"},map[string]string{"url":"/note/golang", "method": "Post"}))
	return  nil
}


type Ram struct {
	DefaultPolicy	[]string
}

func (ram *Ram) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cons := make(map[string]string)
	cons["methd"] = r.Method
	cons["url"] = r.URL.Path
	cons["addr"] = GetRealClientIP(r)
	// read token
	t, err := r.Cookie("t")
	var ok bool
	if err == nil {
		cl := &token.MapClaims{}
		to, err := token.ParseWithClaims(t.Value, cl, func(to *token.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := to.Method.(*token.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", to.Header["alg"])
			}
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte("secret"), nil
		})
		if !to.Valid || err != nil {
			log.Error(err)
		}else {
			ok = Match(strings.Split((*cl)["policy"].(string), ","), cons)
		}
	}
	ok = ok || Match(ram.DefaultPolicy, cons)
	if !ok {
		log.Info("Ram Deny")
		http.Error(w, http.StatusText(403), http.StatusForbidden)
		r.Method = "Deny"
	}
}