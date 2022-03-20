package cuslog

import (
	"bytes"
	"runtime"
	"time"
)

type Entry struct {
	logger *logger
	Buffer *bytes.Buffer
	Map map[string]interface{}
	level Level
	Time time.Time
	File string
	Line int
	Func string
	Format  string
	Args []interface{}
}

func entry(logger *logger) *Entry {
	return &Entry{
		logger: logger,
		Buffer: new(bytes.Buffer),
		Map: make(map[string]interface{}, 5),
	}
}

func (e *Entry) write(level Level, format string, args ...interface{})  {
	if e.logger.opt.level > level {
		return
	}

	e.Format = format
	e.Time = time.Now()
	e.Args = args
	e.level = level

	if !e.logger.opt.disableCaller {
		//if pc, file, line, ok := runtime.Caller(2); !ok {
		if pc, file, line, ok := runtime.Caller(0); !ok {
			e.File = "???"
			e.Func = "???"
		} else {
			e.File, e.Line, e.Func = file, line, runtime.FuncForPC(pc).Name()
			//e.Func = e.Func[strings.LastIndex(e.Func, "/")+1:]
		}
	}

	e.format()
	e.writer()
	e.release()
}

func (e *Entry) format()  {
	_ = e.logger.opt.formatter.Format(e)
}

func (e *Entry) writer()  {
	// 加锁是为了并发安全
	e.logger.mu.Lock()
	_, _ = e.logger.opt.output.Write(e.Buffer.Bytes())
	e.logger.mu.Unlock()
}

func (e *Entry) release()  {
	e.File, e.Line, e.Format, e.Func, e.Args = "", 0, "", "", nil
	e.Buffer.Reset()
	e.logger.entryPool.Put(e)
}
