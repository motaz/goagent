package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func CDRConf(w http.ResponseWriter, r *http.Request) {
	var err string
	type JSONRequest struct {
		Server string
		Duname string
		Dpass  string
		Dname  string
		Ctab   string
		Ckey   string
	}

	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}
	dpath, _ := execShell("locate libmyodbc.so")
	spath, _ := execShell("locate libodbcmyS.so")
	result := JSONResult{true, 0, "", ""}

	w.Header().Add("Content-Type", "application/json")

	body, _ := ioutil.ReadAll(r.Body)
	var jrequest JSONRequest
	er := json.Unmarshal(body, &jrequest)
	if er != nil {
		result.Success = false
		result.Errorcode = 1
		result.Message = er.Error()
	} else {
		//set cdr configuration in configuration file
		if strings.Compare(dpath, "") != 0 && strings.Compare(spath, "") != 0 {
			dpath = dpath[:len(dpath)-1]
			spath = spath[:len(spath)-1]
			msqlcont := "Description     = MySQL driver\nDriver          = " +
				dpath + " \nSetup           = " + spath
			err = addConfNode("/etc/odbcinst.ini", "[Default]", "Driver          = "+dpath)
			addConfNode("/etc/odbcinst.ini", "[MySQL]", msqlcont)
			ser := jrequest.Server
			duname := jrequest.Duname
			dpass := jrequest.Dpass
			dname := jrequest.Dname
			cdrtab := jrequest.Ctab
			cdrkey := jrequest.Ckey
			astdsncont := "Driver          = MySQL\nDescription     = MySQL Connector for Asterisk\nServer          = " +
				ser + "\nPort            = 3306\nDatabase        = " +
				dname + "\nusername        = " +
				duname + "\npassword        = " +
				dpass + "\nOption          = 3\nSocket          = /var/run/mysqld/mysqld.sock"
			err = addConfNode("/etc/odbc.ini", "[MySQL-asterisk]", astdsncont)
			backupFile("res_odbc.conf")
			rsodcont := "enabled=yes\ndsn=MySQL-asterisk\nusername=" +
				duname + "\npassword=" +
				dpass + "\npooling=no\nlimit=1\npre-connect=yes\nshare_connections=yes\nsanitysql=select 1\nisolation=repeatable_read"
			err = addConfNode("/etc/asterisk/res_odbc.conf", "[asterisk]", rsodcont)
			backupFile("cdr_odbc.conf")
			cdrodbccont := "dsn=asterisk\nloguniqueid=yes\ntable=" +
				cdrtab + "\ndispositionstring=yes\nusegmtime=no\nhrtime=yes"
			err = addConfNode("/etc/asterisk/cdr_odbc.conf", "[global]", cdrodbccont)
			backupFile("cdr_manager.conf")
			err = modifyConfNode("/etc/asterisk/cdr_manager.conf", "[general]", "[general]\nenabled = yes")
			backupFile("cdr_adaptive_odbc.conf")
			cdradcont := "connection=asterisk\ntable=" +
				cdrtab + "\nalias start=" +
				cdrkey
			err = addConfNode("/etc/asterisk/cdr_adaptive_odbc.conf", "[asteriskcdr]", cdradcont)
			if strings.Compare(err, "") != 0 {
				result.Success = false
				result.Errorcode = 2
				result.Message = err
				writeLog("Error in CDRConf: " + err)
			} else {
				setConfigParameter("/etc/simpletrunk/stagent.ini", "cdrdbserver", ser)
				setConfigParameter("/etc/simpletrunk/stagent.ini", "cdruser", duname)
				setConfigParameter("/etc/simpletrunk/stagent.ini", "cdrpass", dpass)
				setConfigParameter("/etc/simpletrunk/stagent.ini", "cdrdatabase", dname)
				setConfigParameter("/etc/simpletrunk/stagent.ini", "cdrtable", cdrtab)
				setConfigParameter("/etc/simpletrunk/stagent.ini", "cdrkeyfield", cdrkey)
				execCLI("core reload", r.RemoteAddr)
			}
		} else {

			result.Success = false
			result.Errorcode = 3
			result.Message = "libmyodbc not installed to install do command: sudo apt-get install libmyodbc "
			writeLog("Error in CDRConf: " + result.Message)
		}

	}

	output, _ := json.Marshal(result)

	w.Write(output)

}
func ModifyCDRConf(w http.ResponseWriter, r *http.Request) {
	var err string
	type JSONRequest struct {
		Server string
		Duname string
		Dpass  string
		Dname  string
		Ctab   string
		Ckey   string
	}

	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}
	result := JSONResult{true, 0, "", ""}

	w.Header().Add("Content-Type", "text/html")

	body, _ := ioutil.ReadAll(r.Body)
	var jrequest JSONRequest
	er := json.Unmarshal(body, &jrequest)
	if er != nil {
		result.Success = false
		result.Errorcode = 1
		result.Message = er.Error()
	} else {
		//set cdr configuration in configuration file
		ser := jrequest.Server
		duname := jrequest.Duname
		dpass := jrequest.Dpass
		dname := jrequest.Dname
		cdrtab := jrequest.Ctab
		cdrkey := jrequest.Ckey
		astdsncont := "Driver=MySQL\nDescription=MySQL Connector for Asterisk\nServer          = " +
			ser + "\nPort=3306\nDatabase=" +
			dname + "\nusername=" +
			duname + "\npassword=" +
			dpass + "\nOption=3\nSocket=/var/run/mysqld/mysqld.sock"
		err = modifyConfNode("/etc/odbc.ini", "[MySQL-asterisk]", "[MySQL-asterisk]\n"+astdsncont)
		backupFile("res_odbc.conf")
		rsodcont := "enabled=yes\ndsn=MySQL-asterisk\nusername=" +
			duname + "\npassword=" +
			dpass + "\npooling=no\nlimit=1\npre-connect=yes\nshare_connections=yes\nsanitysql=select 1\nisolation=repeatable_read"
		err = modifyConfNode("/etc/asterisk/res_odbc.conf", "[asterisk]", "[asterisk]\n"+rsodcont)
		backupFile("cdr_odbc.conf")
		cdrodbccont := "dsn=asterisk\nloguniqueid=yes\ntable=" +
			cdrtab + "\ndispositionstring=yes\nusegmtime=no\nhrtime=yes"
		err = modifyConfNode("/etc/asterisk/cdr_odbc.conf", "[global]", "[global]\n"+cdrodbccont)
		cdradcont := "connection=asterisk\ntable=" +
			cdrtab + "\nalias start=" +
			cdrkey
		err = modifyConfNode("/etc/asterisk/cdr_adaptive_odbc.conf", "[asteriskcdr]", "[asteriskcdr]\n"+cdradcont)
		if strings.Compare(err, "") != 0 {
			result.Success = false
			result.Errorcode = 2
			result.Message = err
			writeLog("Error in CDRConf: " + err)
		} else {
			setConfigParameter("/etc/simpletrunk/stagent.ini", "cdrdbserver", ser)
			setConfigParameter("/etc/simpletrunk/stagent.ini", "cdruser", duname)
			setConfigParameter("/etc/simpletrunk/stagent.ini", "cdrpass", dpass)
			setConfigParameter("/etc/simpletrunk/stagent.ini", "cdrdatabase", dname)
			setConfigParameter("/etc/simpletrunk/stagent.ini", "cdrtable", cdrtab)
			setConfigParameter("/etc/simpletrunk/stagent.ini", "cdrkeyfield", cdrkey)
			execCLI("core reload", r.RemoteAddr)
		}
	}
	output, _ := json.Marshal(result)
	w.Write(output)
}
func GetCDRConf(w http.ResponseWriter, r *http.Request) {

	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}
	result := JSONResult{true, 0, "", ""}
	keyf := getConfigValueLocal("cdrkeyfield")
	if keyf == "" {
		keyf = "calldate"
	}
	w.Header().Add("Content-Type", "text/html")
	result.Result = getConfigValueLocal("cdrdbserver") + ":" +
		getConfigValueLocal("cdruser") + ":" +
		getConfigValueLocal("cdrpass") + ":" +
		getConfigValueLocal("cdrdatabase") + ":" +
		getConfigValueLocal("cdrtable") + ":" +
		keyf

	output, _ := json.Marshal(result)
	w.Write(output)
}
func IsCDRConf(w http.ResponseWriter, r *http.Request) {
	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}
	result := JSONResult{true, 0, "", ""}
	w.Header().Add("Content-Type", "text/html")
	if strings.Compare(getConfigValueLocal("cdrdbserver"), "") == 0 {
		result.Success = false
		result.Errorcode = 1
		result.Message = "Not Config yet"
	}
	output, _ := json.Marshal(result)
	w.Write(output)
}
func GetCDRConfStatus(w http.ResponseWriter, r *http.Request) {
	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}
	result := JSONResult{true, 0, "", ""}
	w.Header().Add("Content-Type", "text/html")
	rs, _ := execCLI("odbc show all", r.RemoteAddr)
	if !strings.Contains(rs, "Connected: Yes") {
		result.Success = false
		result.Errorcode = 1
		result.Message = "Configuration Error"
	}
	output, _ := json.Marshal(result)
	w.Write(output)
}
