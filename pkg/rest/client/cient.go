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
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Body       []byte `jsom:"Body"`
}

func (r *Response) String() string {
	msg := r.Message
	if msg == "" {
		msg = string(r.Body)
	}
	return fmt.Sprintf("HTTP response - Status code: %d. Message: '%s'", r.StatusCode, msg)
}

func Get(url string, responseBody interface{}) (*Response, error) {
	return PerformRequest("GET", url, nil, responseBody, http.StatusOK)
}

func Post(url string, requestBody interface{}, responseBody interface{}) (*Response, error) {
	return PerformRequest("POST", url, requestBody, responseBody, http.StatusCreated)
}

func Put(url string, requestBody interface{}, responseBody interface{}) (*Response, error) {
	return PerformRequest("PUT", url, requestBody, responseBody, http.StatusOK)
}
func Delete(url string) (*Response, error) {
	return PerformRequest("DELETE", url, nil, nil, http.StatusOK)
}

func PerformRequest(method string, url string, requestBody interface{}, responseBody interface{}, expectedStatusCode int) (*Response, error) {
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
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error performing %s request equest to [%s]: %s", method, url, err)
	}
	defer res.Body.Close()
	// We need to read all because else keep-alive won't work
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body of %s request to [%s]: %s", method, url, err)
	}
	if expectedStatusCode == 0 {
		expectedStatusCode = http.StatusOK
	}
	if res.StatusCode != expectedStatusCode {
		var response Response
		fmt.Printf("Response: [%s]", data)
		json.Unmarshal(data, &response)
		// We're ignoring the error because the server might have return
		// an unparsable error message, e.g. in HTML.
		// We will store that message in the Body attribute and just go ahead
		response.StatusCode = res.StatusCode
		response.Body = data
		return &response, fmt.Errorf("Unexpected status code: %d", res.StatusCode)
	}
	if responseBody != nil {
		err = json.Unmarshal(data, responseBody)
		if err != nil {
			return &Response{StatusCode: res.StatusCode, Body: data},
				fmt.Errorf("Error parsing response body of %s request to [%s]: %s", method, url, err)
		}
	}
	return &Response{StatusCode: res.StatusCode, Body: data}, nil
}
