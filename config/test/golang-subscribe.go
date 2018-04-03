
package main;

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"io/ioutil"
	"net/http"
	"encoding/base64"
)

var godaemon = flag.Bool("d", false, "run app as a daemon with -d=true or -d true.")

func init() {
	if !flag.Parsed() {
		flag.Parse()
	}

	if *godaemon {
		cmd := exec.Command(os.Args[0], flag.Args()[1:]...)
		cmd.Start()
		fmt.Printf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid)
		*godaemon = false
		os.Exit(0)
	}
}


func main(){
	http.HandleFunc("/",hello)
	http.ListenAndServe(":8001", nil);
}

func hello(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("data%s.txt",r.URL.Path)
	if !PathExist(path){
		return;
	}
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte(base64.RawURLEncoding.EncodeToString(data) ))
}

func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}