package main;

import (  
    "fmt" 
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)  

func main() {
	check("root","");
}

func check(user,pass string) (int,string,int) {
	db, err := sql.Open("mysql","root:@/Jass")
	checkErr(err)
	defer db.Close()
	//rows, err := db.Query("call pro_auth_login(?,?,?,@out_uid,@out_name,@out_level)","asdmin","",1220)
	//checkErr(err)
	var uid = 0
	var name string
	var level int = -1
	db.QueryRow("call pro_auth_login(?,?,?,@out_uid,@out_name,@out_level)",user,pass,0).Scan(&uid,&name,&level)
	//for rows.Next() {
	//	err = rows.Scan(&uid,&name,&level)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(name)
		fmt.Println(level)
	//}
	return uid,name,level
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}


//		db.QueryRow("call pro_auth_login(?,?,?,@out_uid,@out_name,@out_level)",user,pass,0).Scan(&t,&val)