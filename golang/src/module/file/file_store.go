package file;

import (
	"time"
	"net/http"
	"database/sql"
)


type Manager struct {
	storetype 	map[int]Store
	storeinfo 	map[int]*StoreInfo
	storetoken	map[string]Token
}

type StoreInfo struct {
	Uptype 	int
	Uphost 	string
	Key 	string
	Secret 	string
}

type Token struct {
	StoreInfo
	Path  	string
	Size 	int
	Expire 	time.Time
}


type Store interface {
	Policy(calluri,dir string,len int) []byte
	Save(w http.ResponseWriter, r *http.Request)
	Load(path string)
}

func (m *Manager) Load() error {
	if db, err := sql.Open("mysql","root:@/Jass");err==nil {
		defer db.Close()
		var id int
		var uptype int
		var host string
		var key string
		var secret string
		rows, err := db.Query("SELECT ID,Uptype,Uphost,Key,Secret FROM tb_file_store;")
		if err == nil {
			ns := make(map[int]*StoreInfo)
			for rows.Next(){
				rows.Scan(&id,&uptype,&host,&key,&secret)
				ns[id]=&StoreInfo{	Uptype: uptype,	Uphost: host,	Key: key,	Secret: secret 	}
			}
			m.storeinfo = ns
			return nil
		}
		return err
	}
	return nil
}



