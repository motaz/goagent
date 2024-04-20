package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/motaz/codeutils"
)

func getIsWorkingTime(w http.ResponseWriter, r *http.Request) {

	schedule := r.FormValue("schedule")
	if schedule != "" {
		result := isWorkingTime("/etc/asterisk/" + schedule + ".conf")
		res := "0"
		if result {
			res = "1"
		}
		w.Write([]byte(res))
	}

}

func isWorkingTime(file string) (result bool) {

	result = false
	now := time.Now()
	aday := int(now.Weekday()) + 1
	data := codeutils.GetConfigValue(file, strconv.Itoa(aday))
	if data != "" {
		data = strings.TrimSpace(data)
		from := data[:strings.Index(data, " ")]
		from = strings.TrimSpace(from)
		to := data[strings.Index(data, "to")+2:]
		to = strings.TrimSpace(to)

		fromTime, _ := time.Parse("15:04", from)
		toTime, _ := time.Parse("15:04", to)
		ftime := time.Date(now.Year(), now.Month(), now.Day(), fromTime.Hour(), fromTime.Minute(), 0, 0, time.Local)
		ttime := time.Date(now.Year(), now.Month(), now.Day(), toTime.Hour(), toTime.Minute(), 59, 999, time.Local)
		if now.After(ftime) && (now.Before(ttime)) {
			result = true
		}
	}
	return

}
