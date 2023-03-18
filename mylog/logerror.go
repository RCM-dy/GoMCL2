package mylog

import (
	"fmt"
	"os"
	"time"
)

func LogError(f *os.File, err error, pkg string) error {
	_, errs := f.WriteString(fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+" "+pkg+"/ERROR]: %s\n", err.Error()))
	return errs
}
