package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type JSONData struct {
	Header []string   `json:"header"`
	Data   [][]string `json:"data"`
}

type JSONResult struct {
	Success   bool     `json:"success"`
	Errorcode int      `json:"errorcode"`
	Message   string   `json:"message"`
	Result    JSONData `json:"result"`
}

func returnError(w http.ResponseWriter, errorMessage string) {

	var jd JSONData
	result := JSONResult{false, 5, errorMessage, jd}
	output, er := json.Marshal(result)
	if er != nil {
		writeLog("Error in Marshal: " + er.Error())
	}
	w.Write(output)
}

func GetLastCDR(w http.ResponseWriter, r *http.Request) {

	cdrdbserver := getConfigValueLocal("cdrdbserver")
	cdrdatabase := getConfigValueLocal("cdrdatabase")
	cdruser := getConfigValueLocal("cdruser")
	cdrpass := getConfigValueLocal("cdrpass")
	cdrtable := getConfigValueLocal("cdrtable")
	cdrkeyfield := getConfigValueLocal("cdrkeyfield")
	if cdrkeyfield == "" {
		cdrkeyfield = "calldate"
	}

	w.Header().Add("Content-Type", "text/html")
	connection := cdruser + ":" + cdrpass + "@tcp(" + cdrdbserver + ":3306)/" + cdrdatabase + "?charset=utf8"

	db, err := sql.Open("mysql", connection)
	if err != nil {
		println(err.Error())
		writeLog("Error in GetLastCDR db connection: " + err.Error())

		returnError(w, err.Error())

	} else {

		query := "select * from " + cdrtable + " order by " + cdrkeyfield +
			" desc limit 50"

		rows, err := db.Query(query)
		if err != nil {
			writeLog("Error in query: " + err.Error())
			returnError(w, err.Error())
		} else {
			columnNames, _ := rows.Columns()
			size := len(columnNames)

			// Header
			var data JSONData
			data.Header = make([]string, size)
			for i := 0; i < size; i++ {
				data.Header[i] = columnNames[i]
			}

			// Data
			data.Data = make([][]string, 0)

			var fields []interface{}
			for i := 0; i < size; i++ {
				fields = append(fields, new(string))
			}
			for rows.Next() {
				slice := make([]string, size)
				rows.Scan(fields...)
				for i := 0; i < size; i++ {
					text := *(fields[i].(*string))
					slice[i] = text
				}
				data.Data = append(data.Data, slice)

			}
			db.Close()
			// result
			result := JSONResult{true, 0, "", data}

			output, _ := json.Marshal(result)
			w.Write(output)
		}

	}

}
