package main;

import (
	"fmt"
	"net/http"
	"time"
)

type customHandler	struct{

}
    
func(cb *customHandler) ServeHTTP( w http.ResponseWriter, r *http.Request ) {
	fmt.Println("customHandler!!");
	w.Write([]byte("customHandler!!"));
}

func main() {
	var server *http.Server = &http.Server{
		Addr:			":8081",
		Handler:		&customHandler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 <<20,
	}
	server.ListenAndServe();
	select {
	}
}