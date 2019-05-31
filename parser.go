package buttons

import (
	"errors"
	"io/ioutil"
	"net/http"
)

func parseRequest(r *http.Request) (data []byte, err error) {

	token := r.URL.Query().Get("token")
	if secret != token {
		return nil, errors.New("Invalid token")
	}

	return ioutil.ReadAll(r.Body)
}
