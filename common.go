package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func execShell(command string) (string, string) {

	var out bytes.Buffer
	var err bytes.Buffer

	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()
	println("Result: ", out.String())

	println("Error:  ", err.String())
	return out.String(), err.String()
}

func execCLIAsAMI(command string, remoteAddress string) (string, string) {

	// Has been changed to AMI Call
	command = "action:command\ncommand:" + command
	result := actualAMICall("", "", command)

	//result.Message = strings.Replace(result.Message, "\n", "\n\r", -1)
	if strings.Contains(result.Message, "Privilege:") {
		result.Message = result.Message[strings.Index(result.Message, "Privilege:")+19:]
	}
	err := ""
	if !result.Success {
		err = result.Message
	}
	return result.Message, err
}

func execCLI(command string, remoteAddress string) (string, string) {
	writeLog(remoteAddress + ", Executing CLI: " + command)
	result, err := execShell("/usr/sbin/asterisk -rx '" + command + "'")
	return result, err
}

func copyFile(src, dst string) string {
	in, err := os.Open(src)
	if err != nil {
		return "source error: " + err.Error()
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return "dest error: " + err.Error()
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	out.Close()
	if err != nil {
		return err.Error()
	}
	return ""
}

func backupFile(sourceFileName string) {

	t := time.Now()
	atime := t.Format("060102_150405")

	// Check backup directory
	path := "/etc/asterisk/backup/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}

	err := copyFile("/etc/asterisk/"+sourceFileName, "/etc/asterisk/backup/"+sourceFileName+"."+atime)
	if err != "" {
		writeLog("Error while copy: " + err)
	}

}

func getConfigValueDefault(name string, defaultValue string) (val string) {

	val = getConfigValueLocal(name)
	if val == "" {
		val = defaultValue
	}
	return
}

func getConfigValueLocal(name string) string {

	val := GetConfigValue("/etc/simpletrunk/stagent.ini", name)
	val = strings.Replace(val, "\r", "", -1)
	return val
}

func getDefaultConfigFileName(filename string) string {

	if !strings.Contains(filename, "/") {
		filename = "/etc/asterisk/" + filename
	}
	return filename

}
