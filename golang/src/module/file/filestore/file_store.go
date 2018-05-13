package filestore;

import (
	"fmt"
	"path"
	"strings"
	"errors"
	"net/http"
    "crypto/md5"
    "encoding/hex"
	"database/sql"
)


var newstore	map[string]func(string) (Store,error)
var storeproject	map[string]int
var storeint	map[int]Store		// int as storetype
var storestr	map[string]Store	// str as pro path
var localstore	[]int


type Store interface {
	Policy(*PolicyInfo) []byte
	Save(*http.Request) ([]string, error)
//	Load(path string)
//	Del(path string) error
}


type PolicyInfo struct {
	// Host string `json:"host"`
	Directory string `json:"dir"`
	// Expire int64 `json:"expire,omitempty"`
	Length int64 `json:"-"`
	// Size int `json:"size,omitempty"`
	Method string `json:"-"`
	// OSSAccessKeyId string `json:"OSSAccessKeyId,omitempty"`
	// Signature string `json:"signature,omitempty"`
	// Policy string `json:"policy,omitempty"`
	// Callback string `json:"callback,omitempty"`
}

type FileInfo struct {
	Name 	string
	Dir 	bool
	Size 	string
	ModTime	string
}


func get_size(file_bytes int64) string {
    var i     int
    var units = [6]string{"B", "K", "M", "G", "T", "P"}
    i = 0
    for {
        if file_bytes < 1024 {
            return fmt.Sprintf("%d", file_bytes) + units[i]
        }
        file_bytes = file_bytes >> 10
        i++
    }
}
// type Handle struct {
// 	path 	string
// 	fid 	int
// 	sid 	int
// 	store 	Store
// }

// func GetHandle(path string) (Store,error) {
// 	s,ok := storestr[path]
// 	if ok {
// 		return s,nil
// 	}
// 	var n int
// 		if db, err := sql.Open("mysql","root:@/Jass");err==nil {
// 		defer db.Close()
// 		db.QueryRow("SELECT Store FROM tb_file_project WHERE Path=?;",path).Scan(&n)
// 		return n,nil
// 	}else {
// 		return -1,err
// 	}
// 	s,ok = storeint[n]
// 	if ok {
// 		storestr[path] = s
// 		return s,nil
// 	}
// 	return nil,errors.New("undefined store type")
// }


func Register(name string,store func(string) (Store,error)) {
	newstore[name]=store
}

func NewFileStore(name, config string) (Store,error) {
	fs,ok := newstore[name]
	if ok {
		return fs(config)
	}
	return nil,errors.New("undi filestore")
}


func Reload() error {
	if db, err := sql.Open("mysql","root:@/Jass");err==nil {
		defer db.Close()
		var id int
		var store string
		var config string
		rows, err := db.Query("SELECT ID,Name,Config FROM tb_file_store;")
		if err == nil {
			ns := make(map[int]Store)
			for rows.Next(){
				rows.Scan(&id,&store,&config)
				ns[id],err = NewFileStore(store,config)
			}
			storeint = ns
			return nil
		}
		return err
	}
	return nil
}

func Getstore(path string) (Store,error){
	s,ok := storestr[path]
	if ok {
		return s,nil
	}
	n,err := GetPro(path)
	if err!= nil {
		return nil,err
	}
	s,ok = storeint[n]
	if ok {
		storestr[path] = s
		return s,nil
	}
	return nil,errors.New("undefined store type")
}

func GetPro(path string) (int,error) {
	if db, err := sql.Open("mysql","root:@/Jass");err==nil {
		defer db.Close()
		var n int
		db.QueryRow("SELECT Store FROM tb_file_project WHERE Path=?;",path).Scan(&n)
		return n,nil
	}else {
		return -1,err
	}
}

func List(path string) []FileInfo{
	if db, err := sql.Open("mysql","root:@/Jass");err==nil {
		defer db.Close()
		var name string
		var size int64
		var mt string
		rows, err := db.Query("SELECT Name,Size,ModTime FROM tb_file_save WHERE PHash=?;",PathHash(path[6:]))
		fi := make([]FileInfo, 12)
		fmt.Println("filestore List: ",path[6:])
		if err == nil {
			var i int = 0
			for rows.Next(){
				rows.Scan(&name,&size,&mt)
				fi[i]=FileInfo{
					Name: name,
					Size: get_size(size),
					Dir: false,
					ModTime: mt,
				}
				i++
			}
			return fi
		}
		return nil
	}
	return nil
}

func Add(p string,fs []string) {
	fid,err := GetPro(p) 
	if err!= nil {
		return
	}
	if db, err := sql.Open("mysql","root:@/Jass");err==nil {
		defer db.Close()
		for _,f := range fs {
			if strings.HasPrefix(f,"/file/") {
				f = f[6:]
			}
			if stmt, err := db.Prepare("INSERT tb_file_save(FID,Name,Path,Hash,PHash) VALUES(?,?,?,?,?);");err==nil{
				_, err = stmt.Exec(fid, path.Base(f), f, PathHash(f) ,PathHash(path.Dir(f))) 
				fmt.Println("filestore Add: ",f,err)
			}
		} 
	}
}
	
func Del(fs []string) {

}


func PathHash(path string) string{
    md5Ctx := md5.New()
    md5Ctx.Write([]byte(path))
    return hex.EncodeToString(md5Ctx.Sum(nil))
}
func init() {
	newstore = make(map[string]func(string) (Store,error))
	storeproject = make(map[string]int)
	storeint = make(map[int]Store)
	storestr = make(map[string]Store)
}