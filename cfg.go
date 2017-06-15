// cfg - Yet another config file reader
// License: MIT / X11
// Copyright (c) 2013 by James K. Lawless
// jimbo@radiks.net http://www.radiks.net/~jimbo
// http://www.mailsend-online.com
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

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

// Load adds or updates entries in an existing map with string keys
// and string values using a configuration file.
//
// The filename paramter indicates the configuration file to load ...
// the dest parameter is the map that will be updated.
//
// The configuration file entries should be constructed in key=value
// syntax.  A # symbol at the beginning of a line indicates a comment.
// Blank lines are ignored.
func Load(filename string, dest map[string]string) error {
	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	buff := make([]byte, fi.Size())
	f.Read(buff)

	f.Close()
	str := string(buff)
	if !strings.HasSuffix(str, "\n") {
		str = str + "\n"
	}
	//s2 := re.FindAllString(str, -1)
	s2 := strings.Split(str, "\n")

	for i := 0; i < len(s2); {

		if strings.HasPrefix(s2[i], "#") {
			i++
		} else if strings.Contains(s2[i], "=") {
			key := strings.Trim(s2[i][0:strings.Index(s2[i], "=")], " ")
			val := strings.Trim(s2[i][strings.Index(s2[i], "=")+1:len(s2[i])], " ")

			i++
			dest[key] = val
		} else {
			i++
		}
	}
	return nil
}

func GetConfigValue(configFile, name string) string {

	mymap := make(map[string]string)

	err := Load(configFile, mymap)
	if err == nil {
		return mymap[name]
	} else {
		writeLog("Error in GetConfigValue: " + err.Error())
		return ""
	}

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
