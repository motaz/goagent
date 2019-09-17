package main

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func getCurrentDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir

}

func writeLog(event string) {

	t := time.Now()
	today := t.Day()
	old := false
	dir := getCurrentDir() + string(os.PathSeparator) + "log"
	st, err := os.Stat(dir)
	if (err != nil) && (os.IsNotExist(err)) {
		os.Mkdir(dir, 0777)
	}
	if err == nil {
		if st.ModTime().Year() != st.ModTime().Year() {
			old = true
		}
	}
	 logname := dir + string(os.PathSeparator) + "goagent-" + strconv.Itoa(today) + ".log"
	var f *os.File
	if old {
		f, _ = os.OpenFile(logname, os.O_CREATE|os.O_RDWR, 0666)

	} else {
		f, _ = os.OpenFile(logname, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	}
	f.WriteString(t.String() + ": " + event + "\n")
	f.Close()

}
