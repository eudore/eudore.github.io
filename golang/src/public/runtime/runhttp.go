//	by gracedemo
//	https://github.com/tabalt/gracehttp.git
package runhttp

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"sync"
)



const (
	GRACEFUL_ENVIRON_KEY    = "IS_GRACEFUL"
	GRACEFUL_ENVIRON_STRING = GRACEFUL_ENVIRON_KEY + "=1"

	DEFAULT_READ_TIMEOUT  = 60 * time.Second
	DEFAULT_WRITE_TIMEOUT = DEFAULT_READ_TIMEOUT
)

// refer http.ListenAndServe
func ListenAndServe(addr string, handler http.Handler) error {
	return NewServer(addr, handler, DEFAULT_READ_TIMEOUT, DEFAULT_WRITE_TIMEOUT).ListenAndServe()
}

// refer http.ListenAndServeTLS
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	return NewServer(addr, handler, DEFAULT_READ_TIMEOUT, DEFAULT_WRITE_TIMEOUT).ListenAndServeTLS(certFile, keyFile)
}


//Connection
type Connection struct {
	net.Conn
	listener *Listener

	closed bool
}

func (conn *Connection) Close() error {

	if !conn.closed {
		conn.closed = true
		conn.listener.wg.Done()
	}

	return conn.Conn.Close()
}


//Listener
type Listener struct {
	*net.TCPListener

	wg *sync.WaitGroup
}

func NewListener(tl *net.TCPListener) net.Listener {
	return &Listener{
		TCPListener: tl,

		wg: &sync.WaitGroup{},
	}
}

func (l *Listener) Fd() (uintptr, error) {
	file, err := l.TCPListener.File()
	if err != nil {
		return 0, err
	}
	return file.Fd(), nil
}

func (l *Listener) Accept() (net.Conn, error) {

	tc, err := l.AcceptTCP()
	if err != nil {
		return nil, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(time.Minute)

	l.wg.Add(1)

	conn := &Connection{
		Conn:     tc,
		listener: l,
	}
	return conn, nil
}

func (l *Listener) Wait() {
	l.wg.Wait()
}

//Server

// 支持优雅重启的http服务
type Server struct {
	httpServer *http.Server
	listener   net.Listener

	pidfile 	string
	isGraceful bool
	signalChan chan os.Signal
}

// new server
func NewServer(addr string, handler http.Handler, readTimeout, writeTimeout time.Duration) *Server {
	// 获取环境变量
	isGraceful := false
	if os.Getenv(GRACEFUL_ENVIRON_KEY) != "" {
		isGraceful = true
	}

	pidfile := "/var/run/index.pid"
    file, e := os.OpenFile( pidfile,os.O_WRONLY|os.O_TRUNC|os.O_CREATE,0666,)
	if(e!=nil){
		fmt.Println(e)
	}
	defer file.Close()
	// 写字节到文件中
	byteSlice := []byte(fmt.Sprintf("%d",os.Getpid()))
	file.Write(byteSlice)


	// 实例化Server
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,

			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},

		pidfile: pidfile,
		isGraceful: isGraceful,
		signalChan: make(chan os.Signal),
	}
}

func (srv *Server) ListenAndServe() error {
	addr := srv.httpServer.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := srv.getNetTCPListener(addr)
	if err != nil {
		return err
	}

	srv.listener = NewListener(ln)

	return srv.Serve()
}

func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	addr := srv.httpServer.Addr
	if addr == "" {
		addr = ":https"
	}

	config := &tls.Config{}
	if srv.httpServer.TLSConfig != nil {
		*config = *srv.httpServer.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	ln, err := srv.getNetTCPListener(addr)
	if err != nil {
		return err
	}

	srv.listener = tls.NewListener(NewListener(ln), config)
	return srv.Serve()
}

func (srv *Server) Serve() error {
	// 处理信号
	go srv.handleSignals()

	// 处理HTTP请求
	err := srv.httpServer.Serve(srv.listener)

	// 跳出Serve处理代表 listener 已经close，等待所有已有的连接处理结束
	srv.logf("waiting for connection close...")
	srv.listener.(*Listener).Wait()
	srv.logf("all connection closed, process with pid %d shutting down...", os.Getpid())

	return err
}

func (srv *Server) getNetTCPListener(addr string) (*net.TCPListener, error) {
	var ln net.Listener
	var err error

	if srv.isGraceful {
		file := os.NewFile(3, "")
		ln, err = net.FileListener(file)
		if err != nil {
			err = fmt.Errorf("net.FileListener error: %v", err)
			return nil, err
		}
	} else {
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			err = fmt.Errorf("net.Listen error: %v", err)
			return nil, err
		}
	}
	return ln.(*net.TCPListener), nil
}

func (srv *Server) handleSignals() {
	var sig os.Signal

	signal.Notify(
		srv.signalChan,
		syscall.SIGTERM,
		syscall.SIGUSR2,
	)

	pid := os.Getpid()
	for {
		sig = <-srv.signalChan

		switch sig {

		case syscall.SIGTERM | syscall.SIGKILL:

			srv.logf("pid %d received SIGTERM.", pid)
			srv.logf("graceful shutting down http server...")

			// 关闭老进程的连接
			srv.listener.(*Listener).Close()
			os.Remove(srv.pidfile)
			srv.logf("listener of pid %d closed.", pid)

		case syscall.SIGUSR2:

			srv.logf("pid %d received SIGUSR2.", pid)
			srv.logf("graceful restart http server...")

			err := srv.startNewProcess()
			if err != nil {
				srv.logf("start new process failed: %v, pid %d continue serve.", err, pid)
			} else {
				// 关闭老进程的连接
				srv.listener.(*Listener).Close()
				srv.logf("listener of pid %d closed.", pid)
			}

		default:

		}
	}
}

// 启动子进程执行新程序
func (srv *Server) startNewProcess() error {

	listenerFd, err := srv.listener.(*Listener).Fd()
	if err != nil {
		return fmt.Errorf("failed to get socket file descriptor: %v", err)
	}

	path := os.Args[0]

	// 设置标识优雅重启的环境变量
	environList := []string{}
	for _, value := range os.Environ() {
		if value != GRACEFUL_ENVIRON_STRING {
			environList = append(environList, value)
		}
	}
	environList = append(environList, GRACEFUL_ENVIRON_STRING)

	execSpec := &syscall.ProcAttr{
		Env:   environList,
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), listenerFd},
	}

	fork, err := syscall.ForkExec(path, os.Args, execSpec)
	if err != nil {
		return fmt.Errorf("failed to forkexec: %v", err)
	}

	srv.logf("start new process success, pid %d.", fork)

	return nil
}

func (srv *Server) logf(format string, args ...interface{}) {

	if srv.httpServer.ErrorLog != nil {
		srv.httpServer.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}
