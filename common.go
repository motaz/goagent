package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Shell(command string) (string, string) {

	var out bytes.Buffer
	var err bytes.Buffer

	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = &out
	cmd.Stderr = &err
	cmd.Run()
	return out.String(), err.String()
}

func ExecCLI(command string) (string, string) {
	result, err := Shell("/usr/sbin/asterisk -rx '" + command + "'")
	return result, err
}

func CopyFile(src, dst string) string {
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
	var atime string = t.Format("060102_150405")

	// Check backup directory
	path := "/etc/asterisk/backup/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}

	err := CopyFile("/etc/asterisk/"+sourceFileName, "/etc/asterisk/backup/"+sourceFileName+"."+atime)
	if err != "" {
		writeLog("Error while copy: " + err)
	}

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
