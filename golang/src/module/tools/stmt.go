package tools

import (
	"public/log"
	"database/sql"
)

func Stmt(db *sql.DB,s string) *sql.Stmt {
	stmt,err := db.Prepare(s)
	if err == nil {
		return stmt
	}
	log.Error(s,err)
	return nil
}