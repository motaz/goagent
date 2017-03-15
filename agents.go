package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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

			_, er = f.WriteString("agent=> " + jrequest.Agentnumber + "," + jrequest.password + "," + jrequest.Agentname + "\n")
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

func removeAgent(w http.ResponseWriter, r *http.Request) {

	type JSONRequest struct {
		Agentnumber string
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

		// Read entire file
		content, err := ioutil.ReadFile(fileName)
		if err != nil {
			writeLog("Error in readingFile: " + err.Error())
		}
		lines := strings.Split(string(content), "\n")
		f, er := os.Create(fileName)

		if er == nil {
			var found bool = false

			// Write contents without desired agent
			for i := 0; i < len(lines); i++ {

				if strings.Contains(strings.Replace(lines[i], " ", "", -1), "agent=>") &&
					strings.Contains(lines[i], jrequest.Agentnumber+",") {
					// do nothing, skip writing
					found = true
				} else {
					f.WriteString(lines[i] + "\n")
				}
				result.Success = found
				if !found {
					result.Errorcode = 1
					result.Message = "Agent " + jrequest.Agentnumber + " not found"
				}
			}

		} else {
			result.Success = false
			result.Errorcode = 2
			result.Message = er.Error()
			writeLog("Error in removeAgent: " + er.Error())
		}
		f.Close()
	}

	output, _ := json.Marshal(result)

	w.Write(output)

}

func isAgentExist(w http.ResponseWriter, r *http.Request) {

	type JSONRequest struct {
		Agentnumber string
	}

	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}

	result := JSONResult{false, 5, "", ""}
	w.Header().Add("Content-Type", "text/html")

	body, _ := ioutil.ReadAll(r.Body)
	var jrequest JSONRequest
	er := json.Unmarshal(body, &jrequest)
	if er != nil {
		result.Success = false
		result.Errorcode = 1
		result.Message = er.Error()
	} else {
		result = JSONResult{false, 1, "Not found", "Agent " + jrequest.Agentnumber + " not found"}

		// write into file
		fileName := "/etc/asterisk/agents.conf"
		backupFile(fileName)

		// Read entire file
		content, err := ioutil.ReadFile(fileName)
		if err != nil {
			writeLog("Error in readingFile: " + err.Error())
		}
		lines := strings.Split(string(content), "\n")

		if er == nil {

			// search
			for i := 0; i < len(lines); i++ {

				if strings.Contains(strings.Replace(lines[i], " ", "", -1), "agent=>") &&
					strings.Contains(lines[i], jrequest.Agentnumber+",") {
					// do nothing, skip writing
					result.Success = true
					result.Errorcode = 0
					result.Message = "Found"
					result.Result = ""
					break
				}

			}

		} else {
			result.Success = false
			result.Errorcode = 2
			result.Message = er.Error()
			writeLog("Error in removeAgent: " + er.Error())
		}
	}

	output, _ := json.Marshal(result)

	w.Write(output)

}
