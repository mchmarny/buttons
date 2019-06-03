package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// requestHandler handles the HTTP request
func requestHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	token := strings.TrimSpace(r.Header.Get("token"))
	if secret != token {
		logger.Printf("Invalid token: " + token)
		writeResp(w, http.StatusForbidden, "Invalid Token")
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Printf("Error getting data: " + err.Error())
		writeResp(w, http.StatusBadRequest, "Invalid Content")
		return
	}

	if len(data) > 0 {
		logger.Println(string(data))
		err = que.push(r.Context(), data)
		if err != nil {
			logger.Printf("Error storing data: " + err.Error())
			writeResp(w, http.StatusBadRequest, "Internal Error")
			return
		}
	}

	writeResp(w, http.StatusOK, "OK")
	return
}

func writeResp(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(msg)
}
