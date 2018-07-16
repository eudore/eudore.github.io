package store

import (
	"fmt"
	"path"
	"strings"
)

func List(path string) []FileInfo{
	fmt.Println("filestore List: ",path)
	rows, err := stmtQueryListFile.Query(PathHash(path), 0, 50)
	if err != nil {
		return nil
	}
	// fileinfos
	fi := make([]FileInfo, 50)
	var i int = 0
	var name string
	var size int64
	var mt string
	for rows.Next(){
		rows.Scan(&name,&size,&mt)
		fi[i]=FileInfo{
			Name: name,
			Size: getsize(size),
			Dir: false,
			ModTime: mt,
		}
		i++
	}
	return fi
}

func Add(p string,fs []string) {
	// get project store id
	fid,err := GetPro(p) 
	if err!= nil {
		return
	}
	for _,f := range fs {
		if strings.HasPrefix(f,"/file/") {
			f = f[6:]
		}
		_, err := stmtInsertAddFile.Exec(fid, path.Base(f), f, PathHash(f) ,PathHash(path.Dir(f))) 
		fmt.Println("filestore Add: ",f,err)
	}
}
	
func Del(fs []string) {
	for _,f := range fs {
		if strings.HasPrefix(f,"/file/") {
			f = f[6:]
		}
		_, err := stmtDeleteDelFile.Exec(PathHash(f))
		fmt.Println("filestore Add: ",f,err)
	} 
}