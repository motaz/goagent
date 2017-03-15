package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

func addAgent(w http.ResponseWriter, r *http.Request) {

	type JSONRequest struct {
		Agentnumber string
		Agentname   string
		password    string
	}

	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}

	result := JSONResult{true, 0, "", ""}

	w.Header().Add("Content-Type", "text/html")

	body, _ := ioutil.ReadAll(r.Body)
	var jrequest JSONRequest
	er := json.Unmarshal(body, &jrequest)
	if er != nil {
		result.Success = false
		result.Errorcode = 1
		result.Message = er.Error()
	} else {

		// write into file
		fileName := "/etc/asterisk/agents.conf"
		backupFile(fileName)

		f, er := os.OpenFile(fileName, os.O_RDWR+os.O_APPEND+os.O_CREATE, 0666)

		if er == nil {

			_, er = f.WriteString("\nagent=> " + jrequest.Agentnumber + "," + jrequest.password + "," + jrequest.Agentname + "\n")
			if er != nil {
				writeLog("Error in addAgent: " + er.Error())
			}
		} else {
			result.Success = false
			result.Errorcode = 2
			result.Message = er.Error()
			writeLog("Error in addAgent: " + er.Error())
		}
		f.Close()
	}

	output, _ := json.Marshal(result)

	w.Write(output)

}
