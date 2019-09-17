package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type jsonResultType struct {
	Success   bool   `json:"success"`
	Errorcode int    `json:"errorcode"`
	Result    string `json:"result"`
	Message   string `json:"message"`
}

func setControlObject(w http.ResponseWriter, r *http.Request) {

	var controlObjectRequest controlObjectType
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeLog("Error in setControlObject: " + err.Error())
	}

	result := jsonResultType{true, 0, "", ""}

	err = json.Unmarshal(body, &controlObjectRequest)
	if err != nil {
		result = setError(400, err.Error())
	} else {

		err = setControlObjectInDB(controlObjectRequest)
		if err != nil {
			result = setError(500, err.Error())
		}

	}
	res, _ := json.Marshal(result)

	w.Write([]byte(res))

}

func removeControlObject(w http.ResponseWriter, r *http.Request) {

	var controlObjectRequest controlObjectType
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeLog("Error in setControlObject: " + err.Error())
	}

	result := jsonResultType{true, 0, "", ""}

	err = json.Unmarshal(body, &controlObjectRequest)
	if err != nil {
		result = setError(400, err.Error())
	} else {

		err = removeControlObjectFromDB(controlObjectRequest.ObjectType, controlObjectRequest.ObjectName)
		if err != nil {
			result = setError(500, err.Error())
		}

	}
	res, _ := json.Marshal(result)

	w.Write([]byte(res))

}

func setError(errorCode int, errorMessage string) (result jsonResultType) {

	writeLog("Error in setControlObject: " + errorMessage)
	result.Success = false
	result.Errorcode = errorCode
	result.Message = errorMessage
	return
}
