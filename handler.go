package buttons

import (
	"encoding/json"
	"io"
	"net/http"
)

// RequestHandler handles the HTTP request
func RequestHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")

	once.Do(func() {
		if err := configInitializer(r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Printf("Error while initializing configuration: %v", err)
			io.WriteString(w, errorJSON)
			return
		}
	})

	data, err := parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Printf("Error while processing data: %v", err)
		io.WriteString(w, errorJSON)
		return
	}

	err = que.push(r.Context(), data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Printf("Error while storing event: %v", err)
		io.WriteString(w, errorJSON)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)

	return
}
