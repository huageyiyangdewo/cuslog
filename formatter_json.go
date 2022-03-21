package cuslog

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"time"
)

type JsonFormatter struct {
	IgnoreBasicFields bool
}

func (j *JsonFormatter) Format(e *Entry) error {
	if !j.IgnoreBasicFields {
		e.Map["level"] = LevelNameMapping[e.level]
		e.Map["time"] = e.Time.Format(time.RFC3339)

		if e.File != "" {
			e.Map["file"] = e.File + ":" + strconv.Itoa(e.Line)
		}

		switch e.Format {
		case FmtEmptySeparate:
			e.Map["message"] = fmt.Sprint(e.Args...)
		default:
			e.Map["message"] = fmt.Sprintf(e.Format, e.Args...)
		}

		return jsoniter.NewEncoder(e.Buffer).Encode(e.Map)
	}

	switch e.Format {
	case FmtEmptySeparate:
		for _, arg := range e.Args {
			if err := jsoniter.NewEncoder(e.Buffer).Encode(arg); err != nil {
				return err
			}
		}
	default:
		e.Buffer.WriteString(fmt.Sprintf(e.Format, e.Args...))
	}

	return nil
}
