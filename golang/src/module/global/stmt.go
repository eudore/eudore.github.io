package global

import (
	"public/log"
	"database/sql"
)

func Stmt(s string) *sql.Stmt {
	stmt,err := Sql.Prepare(s)
	if err == nil {
		return stmt
	}
	log.Error(s,err)
	return nil
}