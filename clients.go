package main

import (
	"bufio"
	"encoding/json"

	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type AMIJSONResult struct {
	Success   bool   `json:"success"`
	Errorcode int    `json:"errorcode"`
	Result    string `json:"result"`
	Message   string `json:"message"`
}

func callURL(url string, session string) (string, string) {

	asession := ""
	result := ""
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if session != "" {
		req.Header.Add("Cookie", session)
	}

	if url == "" {
		return "empty URL", ""
	}
	resp, err := client.Do(req)

	if err == nil {

		body, err2 := ioutil.ReadAll(resp.Body)
		if err2 == nil {
			result = string(body[:])
			asession = resp.Header.Get("Set-Cookie")

		} else {
			writeLog("Error: " + err2.Error())

		}

	}
	return result, asession

}

func CallAMI(w http.ResponseWriter, r *http.Request) {

	type JSONRequest struct {
		Username string
		Secret   string
		Command  string
	}

	result := AMIJSONResult{true, 0, "", ""}

	w.Header().Add("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeLog("Error in CallAMI parameters: " + err.Error())
		result.Success = false
		result.Errorcode = 1
		result.Message = err.Error()
	} else {
		var jsonRequest JSONRequest
		er := json.Unmarshal(body, &jsonRequest)
		if er != nil {
			writeLog("Erorr in CallAMI: " + er.Error())
			result.Success = false
			result.Errorcode = 5
			result.Message = er.Error()
			writeLog("Error in CallAMI: " + result.Message)
		} else {
			result = actualAMICall(jsonRequest.Command, jsonRequest.Username, jsonRequest.Secret)

		}

	}
	output, _ := json.Marshal(result)
	w.Write(output)

}

func getAMIStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Content-Type", "text/html")
	result := AMIJSONResult{true, 0, "", ""}
	manres, err := execCLI("manager show settings", "", "", r.RemoteAddr)
	if err != "" {
		result.Success = false
		result.Errorcode = 1
		result.Message = "CLI Commnad Error"
	}
	if strings.Contains(manres, "Manager (AMI):             Yes") && strings.Contains(manres, "Web Manager (AMI/HTTP):    Yes") {
		result.Result = "ok:ok"
	}
	if !strings.Contains(manres, "Manager (AMI):             Yes") && strings.Contains(manres, "Web Manager (AMI/HTTP):    Yes") {
		result.Result = "notok:ok"
	}
	if strings.Contains(manres, "Manager (AMI):             Yes") && !strings.Contains(manres, "Web Manager (AMI/HTTP):    Yes") {
		result.Result = "ok:notok"
	}
	if !strings.Contains(manres, "Manager (AMI):             Yes") && !strings.Contains(manres, "Web Manager (AMI/HTTP):    Yes") {
		result.Result = "notok:notok"
	}
	httpres, err2 := execCLI("http show status", "", "", r.RemoteAddr)
	if err2 != "" {
		result.Success = false
		result.Errorcode = 1
		result.Message = "CLI Commnad Error"
	}
	if strings.Contains(httpres, "Server Enabled") {
		result.Result += ":ok"
	} else {
		result.Result += ":notok"
	}
	if result.Result != "" {
		result.Success = true
	}
	output, _ := json.Marshal(result)
	w.Write(output)
}

func actualAMICall(acommand, username, secret string) AMIJSONResult {

	if username == "" {
		username = getConfigValueLocal("amiusername")
		secret = getConfigValueLocal("amisecret")
	}

	amiurl := getConfigValueLocal("amiurl")
	if amiurl[len(amiurl)-1] != '/' {
		amiurl = amiurl + "/"
	}
	// execute command
	fullURL := amiurl + "?action=login&username=" + username + "&secret=" + secret
	resultStr, session := callURL(fullURL, "")

	var result AMIJSONResult
	if strings.Contains(resultStr, "Success") {
		acommand = strings.Replace(acommand, ":", "=", -1)
		acommand = strings.Replace(acommand, "\n", "&", -1)
		acommand = strings.Replace(acommand, "\r", "", -1)
		acommand = strings.Replace(acommand, "  ", " ", -1)
		acommand = strings.Replace(acommand, " ", "%20", -1)
		fullURL = amiurl + "?" + acommand
		resultStr, session = callURL(fullURL, session)
		result.Success = true
		result.Errorcode = 5
		result.Message = resultStr

		// Logout
		callURL(amiurl+"?action=logoff", session)
	} else {
		result.Success = false
		result.Message = resultStr
		result.Errorcode = 1
	}
	return result
}

func getAMIUsersinfo(w http.ResponseWriter, r *http.Request) {
	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}
	w.Header().Add("Content-Type", "text/html")
	result := JSONResult{true, 0, "", ""}
	var res string
	users := strings.Split(getUsers(), ":")
	for i := 0; i < len(users)-1; i++ {
		sec, _ := getConfNodeProperty("/etc/asterisk/manager.conf", users[i], "secret")
		if sec == "" {
			sec = "not set"
		}
		read, _ := getConfNodeProperty("/etc/asterisk/manager.conf", users[i], "read")
		if read == "" {
			read = "not set"
		}
		write, _ := getConfNodeProperty("/etc/asterisk/manager.conf", users[i], "write")
		if write == "" {
			write = "not set"
		}
		sec = strings.TrimSpace(sec)
		read = strings.TrimSpace(read)
		write = strings.TrimSpace(write)
		res += users[i] + ":" + sec + ":" + read + ":" + write + ";"
	}
	result.Result = res
	output, _ := json.Marshal(result)
	w.Write(output)
}

func getAMIUserInfo(w http.ResponseWriter, r *http.Request) {
	type JSONRequest struct {
		Username string
	}
	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}
	w.Header().Add("Content-Type", "text/html")
	result := JSONResult{true, 0, "", ""}
	var res string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeLog("Error in GetAMIInfo: " + err.Error())
		result.Success = false
		result.Errorcode = 1
		result.Message = err.Error()
	} else {
		var jsonRequest JSONRequest
		er := json.Unmarshal(body, &jsonRequest)
		if er != nil {
			writeLog("Erorr in GetAMIInfo: " + er.Error())
			result.Success = false
			result.Errorcode = 5
			result.Message = er.Error()
		} else {
			sec, _ := getConfNodeProperty("/etc/asterisk/manager.conf", jsonRequest.Username, "secret")
			if sec == "" {
				sec = "secret"
			}
			read, _ := getConfNodeProperty("/etc/asterisk/manager.conf", jsonRequest.Username, "read")
			if read == "" {
				read = "all"
			}
			write, _ := getConfNodeProperty("/etc/asterisk/manager.conf", jsonRequest.Username, "write")
			if write == "" {
				write = "all"
			}
			addi := getAMIAdd(jsonRequest.Username)
			if addi == "" {
				addi = "deny=0.0.0.0/0.0.0.0\npermit=127.0.0.1/255.255.255.0"
			}
			sec = strings.TrimSpace(sec)
			read = strings.TrimSpace(read)
			write = strings.TrimSpace(write)
			addi = strings.TrimSpace(addi)
			res += jsonRequest.Username + ":" + sec + ":" + read + ":" + write + ":" + addi
		}
	}
	result.Result = res
	output, _ := json.Marshal(result)
	w.Write(output)
}

func addAMIUser(w http.ResponseWriter, r *http.Request) {
	type JSONRequest struct {
		Username string
		Secret   string
		Read     string
		Write    string
		Addi     string
	}
	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}
	w.Header().Add("Content-Type", "text/html")
	result := JSONResult{true, 0, "", ""}
	var res string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeLog("Error in AddAMIUser: " + err.Error())
		result.Success = false
		result.Errorcode = 1
		result.Message = err.Error()
	} else {
		var jsonRequest JSONRequest
		er := json.Unmarshal(body, &jsonRequest)
		if er != nil {
			writeLog("Erorr in AddAMIUser: " + er.Error())
			result.Success = false
			result.Errorcode = 5
			result.Message = er.Error()
		} else {
			backupFile("manager.conf")
			user := "[" + jsonRequest.Username + "]"
			cont := "secret=" + jsonRequest.Secret + "\nread=" +
				jsonRequest.Read + "\nwrite=" + jsonRequest.Write + "\n" +
				jsonRequest.Addi
			if !isUserExist(user) {

				err := addConfNode("/etc/asterisk/manager.conf", user, cont)
				if err != "" {
					result.Success = false
					result.Errorcode = 5
					result.Message = err
					writeLog("Error in AddAMIUser: " + result.Message)
				} else {
					result.Success = true
					execCLI("core reload", "", "", r.RemoteAddr)
				}
			} else {
				result.Success = false
				result.Errorcode = 3
				result.Message = "This User is already exist"
			}
		}
	}
	result.Result = res
	output, _ := json.Marshal(result)
	w.Write(output)
}

func modifyAMIUser(w http.ResponseWriter, r *http.Request) {
	type JSONRequest struct {
		Username  string
		NUsername string
		Secret    string
		Read      string
		Write     string
		Addi      string
	}
	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}
	w.Header().Add("Content-Type", "text/html")
	result := JSONResult{true, 0, "", ""}
	var res string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeLog("Error in ModifyAMIUser: " + err.Error())
		result.Success = false
		result.Errorcode = 1
		result.Message = err.Error()
	} else {
		var jsonRequest JSONRequest
		er := json.Unmarshal(body, &jsonRequest)
		if er != nil {
			writeLog("Erorr in ModifyAMIUser: " + er.Error())
			result.Success = false
			result.Errorcode = 5
			result.Message = er.Error()
		} else {
			backupFile("manager.conf")
			user := jsonRequest.Username
			nuser := "[" + jsonRequest.NUsername + "]"
			cont := nuser + "\nsecret=" + jsonRequest.Secret + "\nread=" +
				jsonRequest.Read + "\nwrite=" + jsonRequest.Write + "\n" +
				jsonRequest.Addi
			if !isUserExist(nuser) || (nuser == user) {
				err := modifyConfNode("/etc/asterisk/manager.conf", user, cont)
				if err != "" {
					result.Success = false
					result.Errorcode = 5
					result.Message = err
					writeLog("Error in ModifyAMIUser: " + result.Message)
				} else {
					result.Success = true
					execCLI("core reload", "", "", r.RemoteAddr)
				}
			} else {
				result.Success = false
				result.Errorcode = 3
				result.Message = "This User is already exist"
			}
		}
	}
	result.Result = res
	output, _ := json.Marshal(result)
	w.Write(output)
}

func getUsers() string {
	var res string
	f, err := os.Open("/etc/asterisk/manager.conf")
	if err != nil {
		writeLog("Error in Open File: " + err.Error())
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		li := sc.Text()
		if strings.Contains(li, "[") && !strings.Contains(li, "[general]") && !strings.Contains(li, ";[") {
			res += strings.TrimSpace(li) + ":"
		}
	}
	return res
}

func getAMIAdd(user string) string {
	var res string
	f, err := os.Open("/etc/asterisk/manager.conf")
	if err != nil {
		writeLog("Error in Open File: " + err.Error())
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	toggle := false
	for sc.Scan() {
		li := sc.Text()
		if strings.Contains(li, "[") {
			toggle = false
		}
		if strings.Contains(li, user) {
			toggle = true
		}
		if toggle && !strings.Contains(li, user) && !strings.Contains(li, ";") && !strings.Contains(li, "read") && !strings.Contains(li, "write") && !strings.Contains(li, "secret") {
			li = strings.TrimSpace(li)
			res += li + "\n"
		}

	}
	return res
}

func isUserExist(user string) bool {
	var res bool
	f, err := os.Open("/etc/asterisk/manager.conf")
	if err != nil {
		writeLog("Error in Open File: " + err.Error())
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		li := sc.Text()
		li = strings.TrimSpace(li)
		if strings.Contains(li, user) && !strings.HasPrefix(li, ";") {
			res = true
		}
	}
	return res
}
