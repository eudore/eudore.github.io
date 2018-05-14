package server;


import (
	"fmt"
	"os"
	"errors"
	"os/exec"
	"strconv"
	"io/ioutil"
	"syscall"
)

const (
	DAEMON_ENVIRON_KEY		= "IS_DEAMON"
)

type Callback func() error

var reload func() = func() {}

type Cmd struct {
	Cmd 	string
	Pidfile	string
	Pid		int
	Handler Callback
}

func SetReload(call func()) {
	reload = call
}

func Resolve(cmd ,pidfile string,handler Callback ) error {
	c := Cmd{	
		Cmd: cmd,
		Pidfile: pidfile,
		Handler: handler,
	}
	c.readpid()
	var err error
	switch c.Cmd {
	case "start":
		err = c.Start()
	case "daemon":
		err = c.Daemon()
	case "status":
		err = c.Status()
		fmt.Println("status is", err==nil)
	case "stop":
		err = c.Stop()
		fmt.Println("stop is", err==nil)
	case "restart":
		err = c.Restart()
		fmt.Println("restart is", err==nil)
	default:
		err = errors.New("undefined command " + c.Cmd)
		fmt.Println("undefined command ",c.Cmd)
	}
	return err
}


func (c *Cmd) Start() error{
	c.writepid()
	return c.Handler()
}

func (c *Cmd) Daemon() error {
	if os.Getenv(DAEMON_ENVIRON_KEY)=="" {
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(os.Environ(),"IS_DEAMON=1")
		cmd.Start()
		fmt.Printf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid)
		os.Exit(0)
		return nil
	}else{
		return c.Start()
	}
}

func (c *Cmd) Status() error {
	return syscall.Kill(c.Pid,0)
}

func (c *Cmd) Stop() error {
	return syscall.Kill(c.Pid,syscall.SIGTERM)
}

func (c *Cmd) Reload() error {
	return syscall.Kill(c.Pid,syscall.SIGUSR1)
}

func (c *Cmd) Restart() error {
	return syscall.Kill(c.Pid,syscall.SIGUSR2)
}



func (c *Cmd) readpid() int {
	file,err := os.Open(c.Pidfile)  
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()  
	id,err := ioutil.ReadAll(file)
	c.Pid,_ =strconv.Atoi(string(id))
	return c.Pid
}

func (c *Cmd) writepid() {
	file, _ := os.OpenFile( c.Pidfile,os.O_WRONLY|os.O_TRUNC|os.O_CREATE,0666,)
	defer file.Close()
	byteSlice := []byte(fmt.Sprintf("%d",os.Getpid()))
	file.Write(byteSlice)
}