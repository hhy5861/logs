package logs

import (
	"runtime"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_stacktracePool = sync.Pool{
		New: func() interface{} {
			return newProgramCounters(64)
		},
	}

	_stringsBuilderPool = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
)

const stacktraceKey = "stacktrace"

func stackSkip(skip int, flag bool) zapcore.Field {
	skip += 1
	if !flag {
		return zap.StackSkip(stacktraceKey, skip)
	}

	return JsonStacktrace(skip)
}

func JsonStacktrace(skip int) zapcore.Field {
	return zap.Object(stacktraceKey, takeStacktrace(skip+1))
}

func takeStacktrace(skip int) stackList {
	programCounters := _stacktracePool.Get().(*programCounters)
	defer _stacktracePool.Put(programCounters)

	var numFrames int
	for {
		numFrames = runtime.Callers(skip+2, programCounters.pcs)
		if numFrames < len(programCounters.pcs) {
			break
		}

		programCounters = newProgramCounters(len(programCounters.pcs) * 2)
	}

	frames := runtime.CallersFrames(programCounters.pcs[:numFrames])

	var slist stackList
	sb := _stringsBuilderPool.Get().(*strings.Builder)
	sb.Reset()

	defer _stringsBuilderPool.Put(sb)

	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		sb.WriteString(frame.File)
		sb.WriteString(":")
		sb.WriteString(strconv.Itoa(frame.Line))

		slist = append(slist, &stackFrame{
			fn:   frame.Function,
			file: sb.String(),
		})

		sb.Reset()
	}

	return slist
}

type programCounters struct {
	pcs []uintptr
}

func newProgramCounters(size int) *programCounters {
	return &programCounters{make([]uintptr, size)}
}

type stackList []*stackFrame

type stackFrame struct {
	fn   string
	file string
}

func (s stackList) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for i, v := range s {
		_ = enc.AddObject(strconv.Itoa(i), v)
	}

	return nil
}

func (s *stackFrame) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("func", s.fn)
	enc.AddString("file", s.file)

	return nil
}
