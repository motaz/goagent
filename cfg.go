package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"io/ioutil"
)

var re *regexp.Regexp
var pat = "[#].*\\n|\\s+\\n|\\S+[=]|.*\n"

func init() {
	re, _ = regexp.Compile(pat)
}

func GetConfigValue(name string) string {

	return getConfigValueLocal(name)

}
func setConfigParameter(configfile string, param string, value string) string {
	var res string
	if _, err := os.Stat(configfile); os.IsNotExist(err) {
		path := strings.Split(configfile, "/")
		dir, file := getParentAndFile(path)
		os.MkdirAll(dir, 0666)
		os.Create(dir + file)
	}
	if checkParameter(configfile, param) {
		nl := "\n"
		upper := getUpper(configfile, param)
		lower := getLower(configfile, param)
		if strings.Compare(lower, "") == 0 {
			nl = ""
		}
		f, er1 := os.OpenFile(configfile, os.O_TRUNC+os.O_WRONLY, 0666)
		if er1 != nil {
			res = er1.Error()
		}
		defer f.Close()
		con := param + "=" + value + nl
		_, er2 := f.WriteString(upper + con + lower + "\n")
		if er2 != nil {
			res = er2.Error()
		}
	} else {
		f, _ := os.OpenFile(configfile, os.O_APPEND+os.O_WRONLY, 0666)
		defer f.Close()
		con := param + "=" + value + "\n"
		_, er := f.WriteString(con)
		if er != nil {
			res = er.Error()
		}
	}
	return res

}

func getUpper(configfile string, param string) string {

	fh, _ := os.Open(configfile)
	defer fh.Close()
	var res string
	fs := bufio.NewScanner(fh)
	for fs.Scan() {
		li := fs.Text()
		inf := strings.Split(li, "=")
		if strings.Compare(param, inf[0]) == 0 {
			break
		} else {
			res += li + "\n"
		}
	}
	return res
}

func getLower(configfile string, name string) string {
	fh, _ := os.Open(configfile)
	defer fh.Close()
	var res string
	var flag bool
	fs := bufio.NewScanner(fh)
	for fs.Scan() {
		li := fs.Text()
		inf := strings.Split(li, "=")
		if strings.Compare(name, inf[0]) == 0 {
			flag = true
			continue
		}
		if flag {
			res += li + "\n"
		}

	}
	if strings.Compare(res, "") != 0 {
		res = res[0 : len(res)-1]
	}
	return res
}

func checkParameter(configfile string, param string) bool {
	fh, _ := os.Open(configfile)
	var res bool
	defer fh.Close()
	fs := bufio.NewScanner(fh)
	for fs.Scan() {
		li := fs.Text()
		inf := strings.Split(li, "=")
		if strings.Compare(inf[0], param) == 0 {
			res = true
		}
	}
	return res

}

func getParentAndFile(path []string) (string, string) {
	var dir string
	var file string
	var v string
	for i := 0; i < len(path); i++ {
		if i != len(path)-1 {
			v += path[i] + "/"
		}
		if i == len(path)-1 {
			file = path[i]
		}
	}
	dir = v
	return dir, file
}

func addConfNode(fpath string, node string, content string) string {
	var e string
	f, er := os.OpenFile(fpath, os.O_RDWR+os.O_APPEND+os.O_CREATE, 0666)
	defer f.Close()
	if er == nil {

		_, er = f.WriteString("\n" + node + "\n")
		f.WriteString(content + "\n")
		if er != nil {
			writeLog("Error in addConfigNode: " + er.Error())
		}
	} else {
		e = er.Error()
	}
	return e
}

func modifyConfNode(fpath string, node string, ncontent string) string {
	var e string
	content, err := ioutil.ReadFile(fpath)
	if err != nil {
		writeLog("Error in modifyNode: " + err.Error())
	}
	lines := strings.Split(string(content), "\n")

	// write into file
	f, er := os.Create(fpath)
	defer f.Close()
	if er != nil {
		e = er.Error()
	}
	started := false
	found := false
	for i := 0; i < len(lines); i++ {

		// write new node contents after finding it's header
		if (!found) && strings.Contains(lines[i], node) {
			started = true
			found = true
			f.WriteString(ncontent + "\n\n")

		} else if started && strings.Contains(lines[i], "]") && strings.Index(lines[i], "[") < 5 {
			started = false
		}

		if !started {
			f.WriteString(lines[i] + "\n")
		}
	}
	return e
}

func getConfNodeProperty(fpath string, node string, prob string) (string, string) {
	var e string
	var res string
	f, err := os.Open(fpath)
	if err != nil {
		e = err.Error()
		writeLog("Error in modifyNode: " + err.Error())
	}
	defer f.Close()
	toggle := false
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		li := strings.TrimSpace(sc.Text())
		if strings.Contains(li, "[") {
			toggle = false
		}
		if li == node {
			toggle = true
		}
		if toggle && strings.Contains(li, prob) {
			spl := strings.Split(li, "=")
			res = spl[1]
		}
	}
	return res, e
}
