// GoAgent
// SimpleTrunk web service to manage Asterisk
// Origional code written by Motaz Abdel Azim
// Start development:    26-Jan-2017
// Last update 		  16.March.2017

package main

import (
	"net/http"
)

func main() {

	writeLog("GoAgent has started..")

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
	http.HandleFunc("/ReplaceFile", replaceFile)

	// Binary File upload/download
	http.HandleFunc("/ReceiveFile", receiveFile)
	http.HandleFunc("/DownloadFile", downloadFile)

	// Agents
	http.HandleFunc("/AddAgent", addAgent)
	http.HandleFunc("/RemoveAgent", removeAgent)
	http.HandleFunc("/IsAgentExist", isAgentExist)
	http.HandleFunc("/BackupFiles", BackupFiles)

	// Databases
	http.HandleFunc("/GetLastCDR", GetLastCDR)

	//CDR
	http.HandleFunc("/SetCDRConf", CDRConf)
	http.HandleFunc("/GetCDRConf", GetCDRConf)
	http.HandleFunc("/IsCDRConf", IsCDRConf)
	http.HandleFunc("/GetCDRConfStatus", GetCDRConfStatus)
	http.HandleFunc("/ModifyCDRConf", ModifyCDRConf)

	//AMI Configuration
	http.HandleFunc("/GetAMIStatus", GetAMIStatus)
	http.HandleFunc("/GetAMIUsersInfo", GetAMIUsersinfo)
	http.HandleFunc("/GetAMIUserInfo", GetAMIUserInfo)
	http.HandleFunc("/AddAMIUser", AddAMIUser)
	http.HandleFunc("/ModifyAMIUser", ModifyAMIUser)

	//Schedule
	http.HandleFunc("/IsWorkingTime", IsWorkingTime)

	// Queue waiting count
	http.HandleFunc("/WaitingCount", WaitingCount)

	//Test
	http.HandleFunc("/Test", test)

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		writeLog("Error in ListenAndServe: " + err.Error())
	}
	writeLog("GoAgent has closed")

}
