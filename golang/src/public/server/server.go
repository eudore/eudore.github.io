//	source	https://github.com/tabalt/gracehttp/blob/master/server.go

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"golang.org/x/net/http2"
)

const (
	GRACEFUL_ENVIRON_KEY    = "IS_GRACEFUL"
	GRACEFUL_ENVIRON_STRING = GRACEFUL_ENVIRON_KEY + "=1"
	GRACEFUL_LISTENER_FD    = 3
)

// HTTP server that supported graceful shutdown or restart
type Server struct {
	httpServer		*http.Server
	listener		net.Listener

	handerr			error
	isGraceful		bool
	signalChan		chan os.Signal
	shutdownChan	chan bool
}

func NewServer(addr string, handler http.Handler, readTimeout, writeTimeout time.Duration) *Server {
	isGraceful := false
	if os.Getenv(GRACEFUL_ENVIRON_KEY) != "" {
		isGraceful = true

	}
		
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,

			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},

		handerr:      nil,
		isGraceful:   isGraceful,
		signalChan:   make(chan os.Signal),
		shutdownChan: make(chan bool),
	}
}

func (srv *Server) ListenAndServe() error {
	addr := srv.httpServer.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := srv.getNetListener(addr)
	if err != nil {
		return err
	}

	srv.listener = ln
	return srv.Serve()
}

func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	if os.Getenv("LISTEN_HTML2") != "" {
		http2.ConfigureServer(srv.httpServer, &http2.Server{})
	}
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

	ln, err := srv.getNetListener(addr)
	if err != nil {
		return err
	}

	srv.listener = tls.NewListener(ln, config)
	return srv.Serve()
}

func (srv *Server) Serve() error {
	go srv.handleSignals()
	go func(){
		srv.handerr = srv.httpServer.Serve(srv.listener)
		srv.shutdownChan <- true
	}()

	out("waiting for connections closed.\n")
	<-srv.shutdownChan
	out("all connections closed.\n")

	return srv.handerr
}

func (srv *Server) getNetListener(addr string) (net.Listener, error) {
	var ln net.Listener
	var err error

	if srv.isGraceful {
		file := os.NewFile(GRACEFUL_LISTENER_FD, "")
		ln, err = net.FileListener(file)
		out("isGraceful",err)
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
	return ln, nil
}

func (srv *Server) handleSignals() {
	var sig os.Signal

	signal.Notify(
		srv.signalChan,
		syscall.SIGTERM,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	for {
		sig = <-srv.signalChan
		switch sig {
		case syscall.SIGTERM:
			out("received SIGTERM, graceful shutting down HTTP server.")
			srv.shutdownHTTPServer()
		case syscall.SIGUSR1:
			out("received SIGUSR1, graceful reloading HTTP server.")
			reload()
		case syscall.SIGUSR2:
			out("received SIGUSR2, graceful restarting HTTP server.")

			if pid, err := srv.startNewProcess(); err != nil {
				out("start new process failed: ", err,", continue serving.", err)
			} else {
				out("start new process successed, the new pid is ", pid)
				srv.shutdownHTTPServer()
			}
		default:
		}
	}
}

func (srv *Server) shutdownHTTPServer() {
	if err := srv.httpServer.Shutdown(context.Background()); err != nil {
		out("HTTP server shutdown error: ", err)
	} else {
		out("HTTP server shutdown success.")
		srv.shutdownChan <- true
	}
}

// start new process to handle HTTP Connection
func (srv *Server) startNewProcess() (uintptr, error) {
	listenerFd, err := srv.getTCPListenerFd()
	if err != nil {
		return 0, fmt.Errorf("failed to get socket file descriptor: %v", err)
	}

	// set graceful restart env flag
	envs := []string{}
	for _, value := range os.Environ() {
		if value != GRACEFUL_ENVIRON_STRING {
			envs = append(envs, value)
		}
	}
	envs = append(envs, GRACEFUL_ENVIRON_STRING)

	execSpec := &syscall.ProcAttr{
		Env:   envs,
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), listenerFd},
	}

	fork, err := syscall.ForkExec(os.Args[0], os.Args, execSpec)
	if err != nil {
		return 0, fmt.Errorf("failed to forkexec: %v", err)
	}

	return uintptr(fork), nil
}

func (srv *Server) getTCPListenerFd() (uintptr, error) {
	file, err := srv.listener.(*net.TCPListener).File()
	if err != nil {
		return 0, err
	}
	return file.Fd(), nil
}
