package wgx

import (
	"encoding/json"
	"fmt"
	"runtime"
)

func here() string {
	_, f, l, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", f, l)
}
func hereWithDepth(depth int) string {
	_, f, l, _ := runtime.Caller(depth)
	return fmt.Sprintf("%s:%d", f, l)
}

func JSONLine(i interface{}) string {
	r, _ := json.Marshal(i)
	return string(r)
}

func JSON(i interface{}) string {
	r, e := json.MarshalIndent(i, "  ", "  ")
	if e != nil {
		panic(e)
	}
	return string(r)
}
