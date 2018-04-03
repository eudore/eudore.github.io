package log;

import (
	"os"
	"io"
	"testing"
	"public/log"
)


func TestPut(t *testing.T) {
	errFile,_:=os.OpenFile("access.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	out := log.New(io.MultiWriter(os.Stderr,errFile), log.LstdFlags | log.Lshortfile)
	out.Info("info---")
	out.Output(log.DEBUG,1,"test DEBUG")
	for range []int{1,2,3}{
		out.Output(log.INFO,1,"test INFO")
	}
	out.Output(log.WARNING,1,"test WARNING")
}