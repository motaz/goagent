// goagent project main.go
package main

import (
	"net/http"
)

func main() {

	writeLog("GoAgent has started..")
	http.HandleFunc("/Command", command)
	http.HandleFunc("/ListFiles", ListFiles)
	http.HandleFunc("/GetFile", GetFile)
	http.HandleFunc("/GetLogTail", getLogTail)
	http.HandleFunc("/ModifyFile", modifyFile)
	http.HandleFunc("/DownloadFile", downloadFile)
	http.HandleFunc("/Shell", executeShell)
	http.HandleFunc("/CallAMI", CallAMI)
	http.HandleFunc("/BackupFiles", BackupFiles)
	http.HandleFunc("/GetLastCDR", GetLastCDR)
	http.HandleFunc("/AddNode", addNode)
	http.HandleFunc("/ModifyNode", modifyNode)
	http.HandleFunc("/ReceiveFile", receiveFile)

	http.HandleFunc("/AddAgent", addAgent)
	http.HandleFunc("/RemoveAgent", removeAgent)
	http.HandleFunc("/IsAgentExist", isAgentExist)

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		writeLog("Error in ListenAndServe: " + err.Error())
	}
	writeLog("GoAgent has closed")

}
