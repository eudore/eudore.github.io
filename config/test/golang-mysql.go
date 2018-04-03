package main

import (
	"fmt"  
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql","root:@/Jass")
	checkErr(err)
	rows, err := db.Query("SELECT User FROM Jass.tb_User limit 1")
	checkErr(err)
	for rows.Next() {
		var user string
		err = rows.Scan(&user)
		checkErr(err)
		fmt.Println(user)
	}
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}