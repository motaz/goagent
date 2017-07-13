package main

import (
	"archive/zip"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ListFiles(w http.ResponseWriter, r *http.Request) {

	type JSONRequest struct {
		Foldername string
	}

	type FilesResult struct {
		Success   bool     `json:"success"`
		Errorcode int      `json:"errorcode"`
		Files     []string `json:"files"`
	}

	w.Header().Add("Content-Type", "text/html")

	body, _ := ioutil.ReadAll(r.Body)
	var jsonRequest JSONRequest
	er := json.Unmarshal(body, &jsonRequest)
	if er != nil {
	}
	folder := "/etc/asterisk/"
	if jsonRequest.Foldername != "" {
		folder = jsonRequest.Foldername
	}
	var result FilesResult
	result.Errorcode = 0
	result.Success = true
	files, _ := filepath.Glob(folder + "*")
	for _, fname := range files {
		filename := fname[len(folder):]
		result.Files = append(result.Files, filename)

	}

	output, _ := json.Marshal(result)
	w.Write(output)

}

func GetFile(w http.ResponseWriter, r *http.Request) {

	type FileRequest struct {
		Filename string `json:"filename"`
	}

	type FilesResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Content   string `json:"content"`
		Message   string `json:"message"`
		FileTime  string `json:"filetime"`
	}
	var result FilesResult

	result.Success = true
	result.Errorcode = 0

	w.Header().Add("Content-Type", "text/html")

	body, er := ioutil.ReadAll(r.Body)
	if er != nil {
		result.Success = false
		result.Errorcode = 1
		result.Message = er.Error()
	}
	var fr FileRequest
	er = json.Unmarshal(body, &fr)
	if er != nil {
		result.Success = false
		result.Errorcode = 1
		result.Message = er.Error()
	}
	filename := fr.Filename
	if !strings.Contains(filename, "/") {
		filename = "/etc/asterisk/" + filename
	}

	// Read file modification time
	info, er := os.Stat(filename)
	if er == nil {
		result.FileTime = info.ModTime().String()
	}

	file, er := os.Open(filename)
	if er != nil {
		result.Success = false
		result.Errorcode = 1
		result.Message = "Error while opening file: " + er.Error()
	} else {
		cont, er := ioutil.ReadAll(file)
		if er != nil {
			result.Success = false
			result.Errorcode = 1
			result.Message = er.Error()
		}
		file.Close()

		result.Content = string(cont)
	}

	output, _ := json.Marshal(result)
	w.Write(output)

}

func modifyFile(w http.ResponseWriter, r *http.Request) {

	type JSONRequest struct {
		Filename string
		Content  string
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

		if jrequest.Filename[0] != '/' {
			backupFile(jrequest.Filename)
		}

		// write into file
		fileName := jrequest.Filename
		if !strings.Contains(fileName, "/") {
			fileName = "/etc/asterisk/" + fileName
			writeLog(fileName)
		}
		f, er := os.Create(fileName)
		if er == nil {

			f.WriteString(jrequest.Content)
		} else {
			result.Success = false
			result.Errorcode = 2
			result.Message = er.Error()
		}
		f.Close()
	}

	output, _ := json.Marshal(result)

	w.Write(output)

}

func downloadFile(w http.ResponseWriter, r *http.Request) {

	type JSONRequest struct {
		Filename    string
		ContentType string
	}

	type JSONResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Result    string `json:"result"`
		Message   string `json:"message"`
	}

	//	result := JSONResult{true, 0, "", ""}
	var jrequest JSONRequest
	body, er := ioutil.ReadAll(r.Body)
	if er != nil {
	}

	er = json.Unmarshal(body, &jrequest)
	if er != nil {
	}

	afilename := jrequest.Filename
	onlyfilename := afilename[strings.LastIndex(afilename, "/")+1 : len(afilename)]

	var namepart string

	// Convert gsm to wav
	if strings.Contains(onlyfilename, ".gsm") {
		namepart = onlyfilename[0:strings.Index(onlyfilename, ".")]
		command := "sox " + afilename + " /tmp/" + namepart + ".wav"
		Shell(command)

		afilename = "/tmp/" + namepart + ".wav"
	}
	actualFileDownload(afilename, jrequest.ContentType, w)

}

func actualFileDownload(afilename, aContentType string, w http.ResponseWriter) {

	file, e := os.Open(afilename)
	fi, e := file.Stat()
	if e != nil {
		w.Header().Add("encoding", "UTF-8")

		w.Write([]byte("Error: " + e.Error()))

	} else {
		onlyfilename := afilename[strings.LastIndex(afilename, "/")+1 : len(afilename)]

		w.Header().Add("Content-Type", aContentType)
		w.Header().Add("Content-Disposition", "attachment;filename="+onlyfilename)
		w.Header().Add("encoding", "UTF-8")
		w.Header().Add("Content-Length", strconv.FormatInt(fi.Size(), 10))

		read := bufio.NewReader(file)

		data := make([]byte, 8096)
		var totalsize = 0
		for {
			numread, err := read.Read(data)
			if (err != nil) && (err == io.EOF) {
				break
			} else if err != nil {
				writeLog("Error in actualFileDownload: " + err.Error())
				break
			}
			totalsize = totalsize + numread
			w.Write(data[:numread])
		}

	}
}

func zipit(source, target, ext string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ext) {
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if info.IsDir() {
				header.Name += "/"
			} else {
				header.Method = zip.Deflate
			}

			writer, err := archive.CreateHeader(header)
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err

	})

	return err
}

func BackupFiles(w http.ResponseWriter, r *http.Request) {

	type JSONRequest struct {
		Directory string
		Ext       string
		Name      string
	}

	type FilesResult struct {
		Success   bool   `json:"success"`
		Errorcode int    `json:"errorcode"`
		Message   string `json:"message"`
	}

	w.Header().Add("Content-Type", "text/html")

	body, _ := ioutil.ReadAll(r.Body)
	var jsonRequest JSONRequest
	er := json.Unmarshal(body, &jsonRequest)
	if er != nil {
	}
	var result FilesResult
	result.Errorcode = 0
	result.Success = true
	compressedName := "/tmp/pbx-backup.zip"

	zipit(jsonRequest.Directory, compressedName, jsonRequest.Ext)

	actualFileDownload(compressedName, "application/zip", w)

}

func addNode(w http.ResponseWriter, r *http.Request) {

	type JSONRequest struct {
		Filename string
		Nodename string
		Content  string
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

		backupFile(jrequest.Filename)

		// write into file
		fileName := jrequest.Filename
		if !strings.Contains(fileName, "/") {
			fileName = "/etc/asterisk/" + fileName
		}

		f, er := os.OpenFile(fileName, os.O_RDWR+os.O_APPEND+os.O_CREATE, 0666)

		if er == nil {

			_, er = f.WriteString("\n" + jrequest.Nodename + "\n")
			f.WriteString(jrequest.Content + "\n")
			if er != nil {
				writeLog("Error in addNode: " + er.Error())
			}
		} else {
			result.Success = false
			result.Errorcode = 2
			result.Message = er.Error()
			writeLog("Error in addNode: " + er.Error())
		}
		f.Close()
	}

	output, _ := json.Marshal(result)

	w.Write(output)

}

func modifyNode(w http.ResponseWriter, r *http.Request) {

	type JSONRequest struct {
		Filename string
		Nodename string
		Content  string
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

		backupFile(jrequest.Filename)

		// read file first
		fileName := getDefaultConfigFileName(jrequest.Filename)
		content, err := ioutil.ReadFile(fileName)
		if err != nil {
			writeLog("Error in modifyNode: " + err.Error())
		}
		lines := strings.Split(string(content), "\n")

		// write into file
		f, er := os.Create(fileName)
		if er != nil {
			result.Success = false
			result.Errorcode = 2
			result.Message = er.Error()
		}

		started := false
		found := false
		for i := 0; i < len(lines); i++ {

			// write new node contents after finding it's header
			if (!found) && strings.Contains(lines[i], jrequest.Nodename) {
				started = true
				found = true
				f.WriteString(jrequest.Content + "\n\n")

			} else if started && strings.Contains(lines[i], "]") && strings.Index(lines[i], "[") < 5 {
				started = false
			}

			if !started {
				f.WriteString(lines[i] + "\n")
			}

		}
		f.Close()
	}

	output, _ := json.Marshal(result)

	w.Write(output)

}

func receiveFile(w http.ResponseWriter, r *http.Request) {

	type JSONRequest struct {
		Filename string
		Dir      string
		Content  []string
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
		writeLog("Error in receiveFile: " + result.Message)
	} else {

		needConversion := (strings.Contains(jrequest.Filename, ".wav")) && (!strings.Contains(jrequest.Dir, "/monitor"))
		toDir := jrequest.Dir
		if needConversion {
			toDir = "/tmp/"
		}

		// write into file
		fileName := toDir + jrequest.Filename
		os.MkdirAll(toDir, 0555)

		f, er := os.Create(fileName)
		if er == nil {

			// Decode parts
			for i := 0; i < len(jrequest.Content); i++ {
				sDec, _ := base64.StdEncoding.DecodeString(jrequest.Content[i])
				f.Write(sDec)
			}
		} else {
			result.Success = false
			result.Errorcode = 2
			result.Message = er.Error()
		}
		f.Close()

		// Convert to GSM
		if needConversion {
			namepart := jrequest.Filename[0:strings.Index(jrequest.Filename, ".")]
			command := "sox " + fileName + " " + jrequest.Dir + namepart + ".gsm"
			Shell(command)
		}
		result.Message = "written"
	}

	output, _ := json.Marshal(result)

	w.Write(output)

}

func removeNode(w http.ResponseWriter, r *http.Request) {

	type JSONRequest struct {
		Filename string
		Nodename string
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

		backupFile(jrequest.Filename)

		// read file first
		fileName := getDefaultConfigFileName(jrequest.Filename)
		content, err := ioutil.ReadFile(fileName)
		if err != nil {
			writeLog("Error in modifyNode: " + err.Error())
		}
		lines := strings.Split(string(content), "\n")

		// write into file
		f, er := os.Create(fileName)
		if er != nil {
			result.Success = false
			result.Errorcode = 2
			result.Message = er.Error()
		}

		started := false
		found := false
		for i := 0; i < len(lines); i++ {

			// Skip node after finding it's header
			if (!found) && strings.Contains(lines[i], jrequest.Nodename) {
				started = true
				found = true
				// Skipping node lines

			} else if started && strings.Contains(lines[i], "]") && strings.Index(lines[i], "[") < 5 {
				started = false
			}

			if !started {
				f.WriteString(lines[i] + "\n")
			}

		}
		f.Close()
	}

	output, _ := json.Marshal(result)

	w.Write(output)

}
