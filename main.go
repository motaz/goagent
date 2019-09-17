// GoAgent
// SimpleTrunk web service to manage Asterisk
// Origional code written by Motaz Abdel Azim
// Start development:    26-Jan-2017
// Last update 		     27.May.2019

package main

import (
	"net/http"
)

func main() {
	version := "1.1.0"
	println("GoAgent version: " + version)
	writeLog("GoAgent version: " + version + " has started..")

	// Commands
	http.HandleFunc("/Command", command)    // CLI command
	http.HandleFunc("/Shell", executeShell) // Linux shell command
	http.HandleFunc("/CallAMI", CallAMI)

	// Nodes
	http.HandleFunc("/AddNode", addNode)
	http.HandleFunc("/ModifyNode", modifyNode)
	http.HandleFunc("/RemoveNode", removeNode)

	// Files
	http.HandleFunc("/ListFiles", ListFiles)
	http.HandleFunc("/GetFile", GetFile)
	http.HandleFunc("/ModifyFile", modifyFile)
	http.HandleFunc("/GetLogTail", getLogTail)

	// Binary File upload/download
	http.HandleFunc("/ReceiveFile", receiveFile)
	http.HandleFunc("/DownloadFile", downloadFile)

	// Agents
	http.HandleFunc("/AddAgent", addAgent)
	http.HandleFunc("/RemoveAgent", removeAgent)
	http.HandleFunc("/IsAgentExist", isAgentExist)
	http.HandleFunc("/BackupFiles", BackupFiles)

	// Databases
	http.HandleFunc("/GetLastCDR", getLastCDR)

	//CDR
	http.HandleFunc("/SetCDRConf", CDRConf)
	http.HandleFunc("/GetCDRConf", GetCDRConf)
	http.HandleFunc("/IsCDRConf", IsCDRConf)
	http.HandleFunc("/GetCDRConfStatus", GetCDRConfStatus)
	http.HandleFunc("/ModifyCDRConf", ModifyCDRConf)

	//AMI Configuration
	http.HandleFunc("/GetAMIStatus", getAMIStatus)
	http.HandleFunc("/GetAMIUsersInfo", getAMIUsersinfo)
	http.HandleFunc("/GetAMIUserInfo", getAMIUserInfo)
	http.HandleFunc("/AddAMIUser", addAMIUser)
	http.HandleFunc("/ModifyAMIUser", modifyAMIUser)

	//Schedule
	http.HandleFunc("/IsWorkingTime", getIsWorkingTime)

	// Queue waiting count
	http.HandleFunc("/WaitingCount", getWaitingCount)

	// Control Objects
	http.HandleFunc("/SetControlObject", setControlObject)
	http.HandleFunc("/RemoveControlObject", removeControlObject)

	//Test
	http.HandleFunc("/Test", test)

	// HTTP server
	port := getConfigValueDefault("port", "9091")

	println("Listening on port: " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		writeLog("Error in ListenAndServe: " + err.Error())
	}
	writeLog("GoAgent has closed")

}
