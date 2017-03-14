package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func callURL(url string, session string) (string, string) {

	var asession string = ""
	var result string = ""
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if session != "" {
		req.Header.Add("Cookie", session)
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

	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}

	result := JSONResult{true, 0, "", ""}
	amiurl := GetConfigValue("/etc/simpletrunk/stagent.ini", "amiurl")
	if amiurl[len(amiurl)-1] != '/' {
		amiurl = amiurl + "/"
	}

	w.Header().Add("Content-Type", "text/html")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeLog("Error in CallAMI: " + err.Error())
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
			// execute command
			fullURL := amiurl + "?action=login&username=" + jsonRequest.Username + "&secret=" + jsonRequest.Secret

			resultStr, session := callURL(fullURL, "")
			if strings.Contains(resultStr, "Success") {
				jsonRequest.Command = strings.Replace(jsonRequest.Command, ":", "=", -1)
				jsonRequest.Command = strings.Replace(jsonRequest.Command, "\n", "&", -1)
				jsonRequest.Command = strings.Replace(jsonRequest.Command, "\r", "", -1)
				jsonRequest.Command = strings.Replace(jsonRequest.Command, "  ", " ", -1)
				jsonRequest.Command = strings.Replace(jsonRequest.Command, " ", "%20", -1)
				fullURL = amiurl + "?" + jsonRequest.Command
				resultStr, session = callURL(fullURL, session)
				result.Success = true
				result.Errorcode = 5
				result.Message = resultStr
			} else {
				result.Success = false
				result.Message = resultStr
				result.Errorcode = 1
			}

		}

	}
	writeLog("Error in callAMI: " + result.Message)
	output, _ := json.Marshal(result)
	w.Write(output)

}
