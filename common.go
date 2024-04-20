package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/motaz/codeutils"
)

func writeLog(event string) {

	codeutils.WriteToLog(event, "goagent")

}

func execShell(command string) (result string, errMessage string) {

	var out bytes.Buffer
	var err bytes.Buffer

	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()
	fmt.Println("out: ", out.String())
	result = out.String()
	if result == "" {
		errMessage = err.String()
	} else if strings.Contains(result, "No such") {
		errMessage = result
	}
	return
}

func execCLIAsAMI(command, amiuser, amipass string, remoteAddress string) (string, string) {

	// Has been changed to AMI Call
	command = "action:command\ncommand:" + command
	result := actualAMICall(command, amiuser, amipass)

	err := ""
	if !result.Success {
		err = result.Message
	}
	return result.Message, err
}

func execCLI(command, amiuser, amipass string, remoteAddress string) (string, string) {

	result, err := execCLIAsAMI(command, amiuser, amipass, remoteAddress)
	if err != "" && strings.Contains(command, "reload") ||
		strings.Contains(command, "manager") || strings.Contains(command, "http") {
		writeLog(remoteAddress + ", Switching to CLI: " + command)

		result, err = execShell("/usr/sbin/asterisk -rx '" + command + "'")
	}

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

	val := codeutils.GetConfigValue("config.ini", name)
	val = strings.Replace(val, "\r", "", -1)
	return val
}

func getDefaultConfigFileName(filename string) string {

	if !strings.Contains(filename, "/") {
		filename = "/etc/asterisk/" + filename
	}
	return filename

}
