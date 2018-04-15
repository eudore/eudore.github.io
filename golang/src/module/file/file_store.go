


type Manager struct {
	storetype 	map[int]Store
	storeinfo 	map[int]Info
	storetoken	map[string]Token
}

type Info struct {
	Uptype 	int
	Uphost 	string
	Key 	string
	Secret 	string
}

type Token struct {
	Info
	Path  	string
	Size 	int
	Expire 	time.time
}


type Store interface {
	Policy(calluri,dir string,len int) []byte
	UpFile(w http.ResponseWriter, r *http.Request)
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
			ns := make(map[int]FileStore)
			for rows.Next(){
				rows.Scan(&id,&uptype,&host,&key,&Secret)
				auth[id]=&FileStore{	Uptype: uptype,	Uphost: host,	Key: key,	Secret: secret 	}
			}
			m.storeinfo = ns
			return nil
		}
		return err
	}
	return err
}



