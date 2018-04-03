package main  

import (
	"flag"
	"github.com/golang/glog"
)

func main() {
	defer glog.Flush()
	flag.Parse()  
	name := "root"
	glog.Info("Testing glog",name,"s")	
}