package store;

import (
	"errors"
	"net/http"
	"database/sql"
	"module/global"
	"public/log"
)

var stmtQueryStoreConfig,stmtQueryPathStore *sql.Stmt
var (
	stmtQueryFileStore		*sql.Stmt
	smtmQueryPathStore		*sql.Stmt
	stmtQueryIsFile			*sql.Stmt
	stmtQueryListFile		*sql.Stmt
	stmtInsertAddFile		*sql.Stmt
	stmtDeleteDelFile		*sql.Stmt
)


var globalDB *sql.DB

var newstore	map[string]func(string) (Store,error)	// new store func
var allstore 	map[int]Store		// ID - Store
var pathid		map[string]int		// path -> ID
var pathstore	map[string]Store	// path -> Store



type Store interface {
	Signed(*Policy) []byte
	Verify(*Policy) error
	Save(*http.Request, string) error
	Load(http.ResponseWriter, string) error
	Callback(*http.Request) (string, error)
//	Save(string, io.io.Reader) error
//	Load(string) (io.io.Reader, error)
//	Del(path string) error
}


type Policy struct {
	Host string `json:"host"`
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


type Config struct {
	ID		int
	Name	string
	Store	string
	Config	string
}

type Signed struct {
	Id 		string
	policy	string
	signature	string
}

type FileInfo struct {
	Name 	string
	Dir 	bool
	Size 	string
	ModTime	string
}


func Register(name string,store func(string) (Store,error)) {
	newstore[name]=store
}

func NewFileStore(name, config string) (Store,error) {
	fn,ok := newstore[name]
	if ok {
		return fn(config)
	}
	return nil,errors.New("undi filestore")
}


func Reload() error {
	globalDB = global.Sql

	stmtQueryFileStore	=	global.Stmt("SELECT ID,Name,Store,Config FROM tb_file_store")
	smtmQueryPathStore 	=	global.Stmt("SELECT Store FROM tb_file_project WHERE Path=?")

	stmtQueryIsFile		=	global.Stmt("SELECT ID FROM `tb_file_save` WHERE `Hash`=?;")
	stmtQueryListFile	=	global.Stmt("SELECT Name,Size,ModTime FROM tb_file_save WHERE PHash=? LIMIT ?,?;")
	stmtInsertAddFile	=	global.Stmt("INSERT tb_file_save(FID,Name,Path,Hash,PHash) VALUES(?,?,?,?,?);")
	stmtDeleteDelFile	=	global.Stmt("DELETE FROM tb_file_save WHERE Hash=?;")

	rows, err := stmtQueryFileStore.Query()
	if err != nil {
		return err
	}
	ns := make(map[int]Store)
	for rows.Next(){
		var c Config
		rows.Scan(&c.ID,&c.Name,&c.Store,&c.Config)
		ns[c.ID],err = NewFileStore(c.Store,c.Config)
		log.Info("init",c.Name,err)
	}
	allstore = ns
	return nil
}

func Getstore(path string) (Store,error){
	s,ok := pathstore[path]
	if ok {
		return s,nil
	}
	n,err := GetPro(path)
	if err!= nil {
		return nil,err
	}
	s,ok = allstore[n]
	if ok {
		pathstore[path] = s
		return s,nil
	}
	return nil,errors.New("undefined store type")
}

func GetPro(path string) (n int,err error) {
	err = smtmQueryPathStore.QueryRow(path).Scan(&n)
	return
}


func init() {
	newstore = make(map[string]func(string) (Store,error))
	allstore = make(map[int]Store)
	pathid = make(map[string]int)
	pathstore = make(map[string]Store)
}
