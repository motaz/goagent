package main

import (
	"encoding/json"

	"io/ioutil"
	"net/http"
)

func convertToAMI(command string) (amiCommand string) {

	amiCommand = "action:command\ncommand:" + command
	return
}

func command(w http.ResponseWriter, r *http.Request) {

	type Command struct {
		Command  string
		Username string
		Secret   string
	}

	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}

	result := JSONResult{true, 0, "", ""}

	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeLog("Error in command: " + err.Error())
		result.Success = false
		result.Errorcode = 1
		result.Message = err.Error()
	} else {
		var c Command
		er := json.Unmarshal(body, &c)
		if er != nil {
			writeLog("Error in command: " + er.Error())
			result.Success = false
			result.Errorcode = 5
			result.Message = er.Error()
		} else {
			resultStr, err := execCLI(c.Command, c.Username, c.Secret, r.RemoteAddr)
			if err != "" {
				result.Success = false
				result.Errorcode = 6
				result.Message = err
			} else {
				result.Errorcode = 0
				result.Success = true
				result.Result = resultStr
			}
		}

	}
	output, _ := json.Marshal(result)
	w.Write(output)

}

func getLogTail(w http.ResponseWriter, r *http.Request) {

	type LogRequest struct {
		File  string
		Lines string
	}

	type CommandResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Content   string `json:"content"`
		Message   string `json:"message"`
	}

	result := CommandResult{true, 0, "", ""}
	w.Header().Add("Content-Type", "text/html")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeLog("Error in getLogTile: " + err.Error())
		result.Success = false
		result.Errorcode = 1
		result.Message = err.Error()
	} else {

		var logR LogRequest
		er := json.Unmarshal(body, &logR)
		if er != nil {
			writeLog("Error in getLogTail: " + er.Error())

		}

		resultStr, err := execShell("tail --lines=" + logR.Lines + " " + logR.File)
		if err == "" {
			result.Content = resultStr
		} else {
			result.Message = err
			result.Success = false
			result.Errorcode = 5
		}
	}
	output, _ := json.Marshal(result)
	w.Write(output)

}

func executeShell(w http.ResponseWriter, r *http.Request) {

	type Command struct {
		Command string
	}

	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}

	result := JSONResult{true, 0, "", ""}

	w.Header().Add("Content-Type", "text/html")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeLog("Error in executeShell: " + err.Error())
		result.Success = false
		result.Errorcode = 1
		result.Message = err.Error()
	} else {
		var c Command
		er := json.Unmarshal(body, &c)
		if er != nil {
			writeLog("Error in executeShell: " + er.Error())
			result.Success = false
			result.Errorcode = 5
			result.Message = er.Error()
		} else {

			resultStr, err := execShell(c.Command)
			if err != "" {
				result.Success = false
				result.Errorcode = 6
				result.Message = err
			}
			result.Result = resultStr
		}

	}
	output, _ := json.Marshal(result)
	w.Write(output)

}
