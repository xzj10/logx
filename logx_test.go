package logx

import (
	"testing"
)

func TestFn(t *testing.T) {
	Log = NewLogx("debug")
	Debug("abc %d  %d", 123, 222)
	Debug(123, 222)
}

// go test -v logx_test.go logx.go fn.go -run TestFn
