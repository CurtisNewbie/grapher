package log

import (
	"fmt"
	"time"

	"github.com/spf13/cast"
)

func Debugf(pat string, args ...any) {
	fmt.Printf("[DEBUG] "+NowStr()+" "+pat+"\n", args...)
}

func NowStr() string {
	return cast.ToString(time.Now().UnixMilli())
}

func Logf(pat string, args ...any) {
	fmt.Printf(NowStr()+" "+pat+"\n", args...)
}
