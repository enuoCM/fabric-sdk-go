package logging

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var doNothing = func() {}
var mu sync.Mutex
var m = make(map[uint64]int)

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func printTrace(id uint64, name, typ string, indent int, startTime *time.Time, duration *time.Duration) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += "\t"
	}
	if startTime != nil {
		fmt.Printf("g[%06d]:%s%s%s[%s]\n", id, indents, typ, name, startTime)
	}
	if duration != nil {
		now :=time.Now()
		fmt.Printf("g[%06d]:%s%s%s[%s - <%s>]\n", id, indents, typ, name, now, duration )
	}
}

func TraceTime(tags ...string) func() {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return doNothing
	}

	id := GetGID()
	fn := runtime.FuncForPC(pc)
	name := strings.TrimPrefix(fn.Name(), "src.cloudminds.com/blockchain/iam/blockchain-service/")
	if len(tags) > 0 {
		name = name + ":" + strings.Join(tags, ".")
	}

	mu.Lock()
	v := m[id]
	m[id] = v + 1
	mu.Unlock()
	now := time.Now()
	printTrace(id, name, "->", v+1, &now, nil)
	return func() {
		s := time.Since(now)
		mu.Lock()
		v := m[id]
		m[id] = v - 1
		mu.Unlock()
		printTrace(id, name, "<-", v, nil, &s)
	}
}
