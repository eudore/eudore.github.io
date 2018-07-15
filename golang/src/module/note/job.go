package note

import (
	"public/log"
	// "module/global"
)


func RenewHash() error {
	rows, err := stmtQueryPathHash.Query()
	if err != nil {
		return err
	}
	var path,hash,phash,nhash,nphash string
	for rows.Next(){
		rows.Scan(&path,&hash,&phash)
		nhash = PathHash(path)
		nphash = PPathHash(path)
		if hash != nhash || phash != nphash {
			log.Info(path)
			_,err := stmtUpdatePathHash.Exec(nhash,nphash,hash)
			if err != nil {
				log.Error(err)
			}
		}
	}
	return nil
}

func UpIndex() error {
	return nil
}