package log;

import (
	"io"
	"os"
	"fmt"
	"sync"
	"time"
	"runtime"
)
 const (
	Ldate			= 1 << iota		// the date in the local time zone: 2009/01/23
	Ltime							// the time in the local time zone: 01:23:23
	Lmicroseconds					// microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile						// full file name and line number: /a/b/c/d.go:23
	Lshortfile						// final file name element and line number: d.go:23. overrides Llongfile
	LUTC							// if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags		= Ldate | Ltime	// initial values for the standard logger
)

const (
	DEBUG		= iota
	INFO
	WARNING
	ERROR
	FATAL
)


type toutput func(level,calldepth int, s string) error

var levelName [5][]byte =  [5][]byte{
	[]byte("DEBUG "),
	[]byte("INFO "),
	[]byte("WARNING "),
	[]byte("ERROR "),
	[]byte("FATAL "),
}


type Logger struct {
	mu 		sync.Mutex
	level	int
	flag	int 
	out		io.Writer
	buf    	[]byte
}


type flushSyncWriter interface {
	Flush() error
	Sync() error
	io.Writer
}

func New(out io.Writer, flag int) *Logger {
	return &Logger{out: out, flag: flag}
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}


func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *Logger) formatHeader(buf *[]byte, level int, file string, line int) {
	//*buf = append(*buf, l.prefix...)
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		t := time.Now()
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	*buf = append(*buf, levelName[level]...)
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}


func (l *Logger) Output(level,calldepth int, s string) error {
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, level, file, line)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}


func (l *Logger) Debug(args ...interface{}) {
	l.Output(DEBUG,2,fmt.Sprint(args...))
} 
func (l *Logger) Info(args ...interface{}) {
	l.Output(INFO,2,fmt.Sprint(args...))
} 
func (l *Logger) Warning(args ...interface{}) {
	l.Output(WARNING,2,fmt.Sprint(args...))
} 
func (l *Logger) Error(args ...interface{}) {
	l.Output(ERROR,2,fmt.Sprint(args...))
} 
func (l *Logger) Fatal(args ...interface{}) {
	l.Output(FATAL,2,fmt.Sprint(args...))
} 
var timeNow = time.Now // Stubbed out for testing.




var out *Logger

func init() {
	errFile,_:=os.OpenFile("access.log",os.O_CREATE|os.O_WRONLY|os.O_APPEND,0666)
	out = New(io.MultiWriter(os.Stderr,errFile), LstdFlags | Lshortfile)
}
func Debug(args ...interface{}) {
	out.Output(DEBUG,2,fmt.Sprint(args...))
} 
func Info(args ...interface{}) {
	out.Output(INFO,2,fmt.Sprint(args...))
} 
func Warning(args ...interface{}) {
	out.Output(WARNING,2,fmt.Sprint(args...))
} 
func Error(args ...interface{}) {
	out.Output(ERROR,2,fmt.Sprint(args...))
} 
func Fatal(args ...interface{}) {
	out.Output(FATAL,2,fmt.Sprint(args...))
}