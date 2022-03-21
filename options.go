package cuslog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	FmtEmptySeparate = ""
)

type Level uint8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

var LevelNameMapping = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel: "INFO",
	WarnLevel: "WARN",
	ErrorLevel: "ERROR",
	PanicLevel: "PANIC",
	FatalLevel: "FATAL",
}

var errUnmarshallLevel = errors.New("can't unmarshall a nil *Level")


// 平常使用时，是通过 xxx.info 这种形式来写日志的，所以函数的参数是 byte?
func (l *Level) unmarshallText(text []byte) bool {
	switch string(text) {
	case "debug", "DEBUG":
		*l = DebugLevel
	case "info", "INFO":
		*l = InfoLevel
	case "warn", "WARN":
		*l = WarnLevel
	case "error", "ERROR":
		*l = ErrorLevel
	case "fatal", "FATAL":
		*l = FatalLevel
	default:
		return false
	}
	return true
}

func (l *Level) UnmarshallText(text []byte) error {
	if l == nil {
		return errUnmarshallLevel
	}

	if l.unmarshallText(text) && l.unmarshallText(bytes.ToLower(text)) {
		// %q 输出的字符串 加上双引号
		return fmt.Errorf("unrecognized level:%q", text)
	}

	return nil
}

type options struct {
	output	io.Writer
	level Level
	stdLevel Level
	formatter Formatter
	disableCaller bool
}

type Option func(*options)

func initOptions(opts ...Option) (o *options) {
	o = &options{}
	for _, opt := range opts {
		opt(o)
	}

	if o.output == nil {
		o.output = os.Stderr
	}

	if o.formatter == nil {
		o.formatter = &TextFormatter{}
	}

	return
}

func WithOutput(output io.Writer) Option {
	return func(o *options) {
		o.output = output
	}
}

func WithLevel(level Level) Option {
	return func(o *options) {
		o.level = level
	}
}

func WithStdLevel(stdLevel Level) Option {
	return func(o *options) {
		o.stdLevel = stdLevel
	}
}

func WithFormatter(formatter Formatter) Option {
	return func(o *options) {
		o.formatter = formatter
	}
}

func WithDisableCaller(caller bool) Option {
	return func(o *options) {
		o.disableCaller = caller
	}
}