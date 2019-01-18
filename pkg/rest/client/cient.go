package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	httpClient = &http.Client{
		Timeout: time.Second * 10,
	}
)

type Response struct {
	StatusCode int
	Message    string
	Body       []byte
}

func (r *Response) String() string {
	msg := r.Message
	if msg == "" {
		msg = string(r.Body)
	}
	return fmt.Sprintf("HTTP response - Status code: %d. Message: '%s'", r.StatusCode, msg)
}

func Get(url string, responseBody interface{}) (*Response, error) {
	return PerformRequest("GET", url, nil, responseBody)
}

func Post(url string, requestBody interface{}, responseBody interface{}) (*Response, error) {
	return PerformRequest("POST", url, requestBody, responseBody)
}

func Put(url string, requestBody interface{}, responseBody interface{}) (*Response, error) {
	return PerformRequest("PUT", url, requestBody, responseBody)
}
func Delete(url string) (*Response, error) {
	return PerformRequest("DELETE", url, nil, nil)
}

func PerformRequest(method string, url string, requestBody interface{}, responseBody interface{}) (*Response, error) {
	if method == "" {
		return nil, fmt.Errorf("No method provided for HTTP request")
	}
	if url == "" {
		return nil, fmt.Errorf("No url provided for HTTP request")
	}
	var reader io.Reader
	if requestBody != nil {
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(requestBody)
		if err != nil {
			return nil, fmt.Errorf("Error marshalling request body for %s request to [%s]: %s", method, url, err)
		}
		reader = buf
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("Error creating %s request equest to [%s]: %s", method, url, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	var statusCode int
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error performing %s request equest to [%s]: %s", method, url, err)
	}
	defer res.Body.Close()
	statusCode = res.StatusCode
	// We need to read all because else keep-alive won't work
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body of %s request to [%s]: %s", method, url, err)
	}
	if statusCode == http.StatusOK || statusCode == http.StatusCreated {
		if responseBody != nil {
			err = json.Unmarshal(data, responseBody)
			if err != nil {
				return &Response{StatusCode: statusCode, Body: data},
					fmt.Errorf("Error parsing response body of %s request to [%s]: %s", method, url, err)
			}
		}
		return &Response{StatusCode: statusCode, Body: data}, nil
	}
	var response Response
	json.Unmarshal(data, &response)
	// We're ignoring the error because the server might have returnes
	// a generic error message. We will store that message in the RawBody
	// and just go ahead
	response.StatusCode = statusCode
	response.Body = data
	return &response, nil
}
