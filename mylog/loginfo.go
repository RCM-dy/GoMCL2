package mylog

import (
	"fmt"
	"os"
	"time"
)

func LogInfo(f *os.File, infos string, pkg string) error {
	_, errs := f.WriteString(fmt.Sprintf("["+time.Now().Format("2006-01-02 15:04:05")+" %s/INFO]: "+infos+Newline, pkg))
	return errs
}
