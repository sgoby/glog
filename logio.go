// User: szh
// Date: 2020/9/24
// Time: 20:00

package glog

import (
	"io"
	"time"
	"runtime"
	"os"
	"sync"
)

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

//
type ILogger interface {
	Output(calldepth int , s string) error
}

//
type Logio struct {
	bufw *Writer
	flag int    // properties
	buf  []byte // for accumulating text to write
	mu   sync.Mutex
}

//
func NewLogio(flag int) *Logio {
	return &Logio{flag: flag,bufw:NewWriter(os.Stdout)}
}

//
func (l *Logio) Flush() error {
	if l.bufw != nil{
		return l.bufw.Flush()
	}
	return nil
}

func (l *Logio) SetWriter(w io.Writer) {
	if l.bufw == nil {
		l.bufw = NewWriter(w)
	}else {
		l.bufw.Flush()
	}
}

//
func (l *Logio) ResetWriter(w io.Writer) {
	if l.bufw != nil{
		l.bufw.ResetWriter(w)
	}
}

//
func (l *Logio) Reset(w io.Writer) {
	if l.bufw != nil{
		l.bufw.Reset(w)
	}
}

//
func (l *Logio) Output(calldepth int , s string) error {
	return l.OutputByLv(calldepth,"",s)
}

//
func (l *Logio) OutputByLv(calldepth int,lv , s string) error {
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	//
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		//l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		//l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line,lv)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.bufw.Write(l.buf)
	return err
}



// formatHeader writes log header to buf in following order:
//   * l.prefix (if it's not blank and Lmsgprefix is unset),
//   * date and/or time (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided),
//   * l.prefix (if it's not blank and Lmsgprefix is set).
func (l *Logio) formatHeader(buf *[]byte, t time.Time, file string, line int,prefix string) {
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
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
	//
	if len(prefix) > 0 {
		*buf = append(*buf, "["...)
		*buf = append(*buf, prefix...)
		*buf = append(*buf, "] "...)
	}
	//
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

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
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