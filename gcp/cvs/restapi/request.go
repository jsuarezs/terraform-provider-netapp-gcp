package restapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2/google"
)

// Request represents a request to a REST API
type Request struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// BuildHTTPReq builds an HTTP request to carry out the REST request
func (r *Request) BuildHTTPReq(host string, serviceAccount string, credentials string, audience string, baseURL string) (*http.Request, error) {
	var keyBytes []byte
	var err error
	var req *http.Request
	url := host + baseURL
	if r.Method != "GET" && r.Method != "DELETE" {
		bodyJSON, err := json.Marshal(r.Params)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(r.Method, url, bytes.NewReader(bodyJSON))
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(r.Method, url, nil)
		if err != nil {
			return nil, err
		}
	}
	if credentials != "" {
		keyBytes = []byte(credentials)
	} else {
		keyBytes, err = ioutil.ReadFile(serviceAccount)
		if err != nil {
			return nil, fmt.Errorf("Unable to read service account key file  %v", err)
		}
	}

	tokenSource, err := google.JWTAccessTokenSourceFromJSON(keyBytes, audience)
	if err != nil {
		return nil, fmt.Errorf("Error building JWT access token source: %v", err)
	}
	jwt, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("Unable to generate JWT token: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwt.AccessToken)

	return req, nil
}
