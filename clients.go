package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"os"
	"bufio"
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
func GetAMIStatus(w http.ResponseWriter , r *http.Request) {
	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}
	w.Header().Add("Content-Type", "text/html")
	result := JSONResult{true, 0, "", ""}
	manres,err:=ExecCLI("manager show settings")
	if err!=""{
		result.Success=false
		result.Errorcode=1
		result.Message="CLI Commnad Error"
	}
	if strings.Contains(manres,"Manager (AMI):             Yes") && strings.Contains(manres,"Web Manager (AMI/HTTP):    Yes"){
		result.Result="ok:ok"
	}
	if !strings.Contains(manres,"Manager (AMI):             Yes") && strings.Contains(manres,"Web Manager (AMI/HTTP):    Yes"){
		result.Result="notok:ok"
	}
	if strings.Contains(manres,"Manager (AMI):             Yes") && !strings.Contains(manres,"Web Manager (AMI/HTTP):    Yes"){
		result.Result="ok:notok"
	}
	if !strings.Contains(manres,"Manager (AMI):             Yes") && !strings.Contains(manres,"Web Manager (AMI/HTTP):    Yes"){
		result.Result="notok:notok"
	}
	httpres,err2:=ExecCLI("http show status")
	if err2!=""{
		result.Success=false
		result.Errorcode=1
		result.Message="CLI Commnad Error"
	}
	if strings.Contains(httpres,"Server Enabled"){
		result.Result+=":ok"
	}else {
		result.Result+=":notok"
	}
	if result.Result!=""{
		result.Success=true
	}
	output, _ := json.Marshal(result)
	w.Write(output)
}
func GetAMIUsersinfo(w http.ResponseWriter , r *http.Request) {
	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}
	w.Header().Add("Content-Type", "text/html")
	result := JSONResult{true, 0, "", ""}
	var res string
	users:=strings.Split(getUsers(),":")
	for i:=0;i<len(users)-1;i++{
		sec,_:=getConfNodeProperty("/etc/asterisk/manager.conf",users[i],"secret")
		if sec==""{
			sec="not set"
		}
		read,_:=getConfNodeProperty("/etc/asterisk/manager.conf",users[i],"read")
		if read==""{
			read="not set"
		}
		write,_:=getConfNodeProperty("/etc/asterisk/manager.conf",users[i],"write")
		if write==""{
			write="not set"
		}
		sec=strings.TrimSpace(sec)
		read=strings.TrimSpace(read)
		write=strings.TrimSpace(write)
		res+=users[i]+":"+sec+":"+read+":"+write+";"
	}
	result.Result=res
	output, _ := json.Marshal(result)
	w.Write(output)
}
func GetAMIUserInfo(w http.ResponseWriter , r *http.Request) {
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
			sec, _ := getConfNodeProperty("/etc/asterisk/manager.conf",jsonRequest.Username, "secret")
			if sec == "" {
				sec = "secret"
			}
			read, _ := getConfNodeProperty("/etc/asterisk/manager.conf",jsonRequest.Username, "read")
			if read == "" {
				read = "all"
			}
			write, _ := getConfNodeProperty("/etc/asterisk/manager.conf", jsonRequest.Username, "write")
			if write == "" {
				write = "all"
			}
			addi:= getAMIAdd(jsonRequest.Username)
			if addi == "" {
				addi = "deny=0.0.0.0/0.0.0.0\npermit=127.0.0.1/255.255.255.0"
			}
			sec = strings.TrimSpace(sec)
			read = strings.TrimSpace(read)
			write = strings.TrimSpace(write)
			addi = strings.TrimSpace(addi)
			res += jsonRequest.Username + ":" + sec + ":" + read + ":" + write+":"+ addi
		}
	}
	result.Result=res
	output, _ := json.Marshal(result)
	w.Write(output)
}
func AddAMIUser(w http.ResponseWriter , r *http.Request) {
	type JSONRequest struct {
		Username string
		Secret string
		Read string
		Write string
		Addi string
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
			writeLog("Error in AddAMIUser: " + result.Message)
		} else {
			user:="["+jsonRequest.Username+"]"
			cont:="secret="+jsonRequest.Secret+"\nread="+
			jsonRequest.Read+"\nwrite="+jsonRequest.Write+"\n"+
			jsonRequest.Addi
			if !isUserExist(user){

				err:=addConfNode("/etc/asterisk/manager.conf",user,cont)
				if err !=""{
					result.Success = false
					result.Errorcode = 5
					result.Message = err
					writeLog("Error in AddAMIUser: " + result.Message)
				}else {
					result.Success = true
					ExecCLI("core reload")
				}
			}else {
				result.Success = false
				result.Errorcode = 3
				result.Message = "This User is already exist"
			}
		}
	}
	result.Result=res
	output, _ := json.Marshal(result)
	w.Write(output)
}
func ModifyAMIUser(w http.ResponseWriter , r *http.Request) {
	type JSONRequest struct {
		Username string
		NUsername string
		Secret string
		Read string
		Write string
		Addi string
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
			writeLog("Error in AddAMIUser: " + result.Message)
		} else {
			user:=jsonRequest.Username
			nuser:="["+jsonRequest.NUsername+"]"
			cont:=nuser+"\nsecret="+jsonRequest.Secret+"\nread="+
				jsonRequest.Read+"\nwrite="+jsonRequest.Write+"\n"+
				jsonRequest.Addi
			if !isUserExist(nuser) || (nuser==user){
				err:=modifyConfNode("/etc/asterisk/manager.conf",user,cont)
				if err !=""{
					result.Success = false
					result.Errorcode = 5
					result.Message = err
					writeLog("Error in AddAMIUser: " + result.Message)
				}else {
					result.Success = true
					ExecCLI("core reload")
				}
			}else {
				result.Success = false
				result.Errorcode = 3
				result.Message = "This User is already exist"
			}
		}
	}
	result.Result=res
	output, _ := json.Marshal(result)
	w.Write(output)
}
func getUsers() string{
	var res string
	f, err := os.Open("/etc/asterisk/manager.conf")
	if err != nil {
		writeLog("Error in modifyNode: " + err.Error())
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		li:=sc.Text()
		if strings.Contains(li,"[")  && !strings.Contains(li,"[general]") && !strings.Contains(li,";["){
			res+=strings.TrimSpace(li)+":"
		}
	}
	return res
}
func getAMIAdd(user string) string{
	var res string
	f, err := os.Open("/etc/asterisk/manager.conf")
	if err != nil {
		writeLog("Error in modifyNode: " + err.Error())
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	toggle:=false
	for sc.Scan() {
		li:=sc.Text()
		if strings.Contains(li,"["){
			toggle=false
		}
		if strings.Contains(li,user) {
			toggle=true
		}
		if toggle && !strings.Contains(li,user) && !strings.Contains(li,";") && !strings.Contains(li,"read") && !strings.Contains(li,"write") && !strings.Contains(li,"secret"){
			li=strings.TrimSpace(li)
			res +=li+"\n"
		}

	}
	return res
}
func isUserExist(user string) bool{
	var res bool
	f, err := os.Open("/etc/asterisk/manager.conf")
	if err != nil {
		writeLog("Error in modifyNode: " + err.Error())
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		li:=sc.Text()
		li=strings.TrimSpace(li)
		if strings.Contains(li,user) && !strings.HasPrefix(li,";") {
			res=true
		}
	}
	return res
}