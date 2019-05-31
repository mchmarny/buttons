package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// requestHandler handles the HTTP request
func requestHandler(w http.ResponseWriter, r *http.Request) {

	token := strings.TrimSpace(r.Header.Get("token"))
	if secret != token {
		logger.Fatal("Invalid token: " + token)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Fatal("Error getting data: " + err.Error())
	}

	if len(data) > 0 {
		logger.Println(string(data))
		err = que.push(r.Context(), data)
		if err != nil {
			logger.Fatal("Error storing data: " + err.Error())
		}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")

	return
}
