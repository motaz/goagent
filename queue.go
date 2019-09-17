package main

import (
	"net/http"
	"strings"
)

func getWaitingCount(w http.ResponseWriter, r *http.Request) {

	queue := r.FormValue("queue")

	result := actualAMICall("", "", "action:command\ncommand:queue show "+queue)
	res := "0"
	if result.Success {
		line := result.Message
		//6120 has 68 calls (max unlimited) in 'rrmemory' strategy (439s holdtime, 113s talktime), W:0, C:2212, A:5083, SL:0.5% within 60s

		line = line[strings.Index(line, "has")+3 : strings.Index(line, "calls")]
		line = strings.Trim(line, " ")
		res = line
	}
	w.Write([]byte(res))

}
