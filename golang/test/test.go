package main;

import (  
	"fmt" 
	"public/server"
	"public/log"
	"public/router"
	"net/http"
	"os"
	"reflect"
)

func main(){
	fmt.Println("------")
	//mux := http.NewServeMux()

	mux := router.Instance()
	mux.GetFunc("/index/ss/", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("hello, golang!\n"))
		log.Info("ssss")
		log.Error("ssss")
	})
	
	fmt.Println(reflect.TypeOf(mux))
	server.Resolve(os.Args[1],"/var/run/test.pid", func() error {
		return server.ListenAndServe(":7070", mux)
	})
}