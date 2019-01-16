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

func (R *Response) String() string {
	msg := R.Message
	if msg == "" {
		msg = string(R.Body)
	}
	return fmt.Sprintf("HTTP response - Status code: %d. Message: '%s'", R.StatusCode, msg)
}

func Get(URL string, ResponseBody interface{}) (*Response, error) {
	return PerformRequest("GET", URL, nil, ResponseBody)
}

func Post(URL string, RequestBody interface{}, ResponseBody interface{}) (*Response, error) {
	return PerformRequest("POST", URL, RequestBody, ResponseBody)
}

func Put(URL string, RequestBody interface{}, ResponseBody interface{}) (*Response, error) {
	return PerformRequest("PUT", URL, RequestBody, ResponseBody)
}
func Delete(URL string) (*Response, error) {
	return PerformRequest("DELETE", URL, nil, nil)
}

func PerformRequest(Method string, URL string, RequestBody interface{}, ResponseBody interface{}) (*Response, error) {
	if Method == "" {
		return nil, fmt.Errorf("No method provided for HTTP request")
	}
	if URL == "" {
		return nil, fmt.Errorf("No URL provided for HTTP request")
	}
	var reader io.Reader
	if RequestBody != nil {
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(RequestBody)
		if err != nil {
			return nil, fmt.Errorf("Error marshalling request body for %s request to [%s]: %s", Method, URL, err)
		}
		reader = buf
	}
	req, err := http.NewRequest(Method, URL, reader)
	if err != nil {
		return nil, fmt.Errorf("Error creating %s request equest to [%s]: %s", Method, URL, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	var statusCode int
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error performing %s request equest to [%s]: %s", Method, URL, err)
	}
	defer res.Body.Close()
	statusCode = res.StatusCode
	// We need to read all because else keep-alive won't work
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body of %s request to [%s]: %s", Method, URL, err)
	}
	if statusCode == http.StatusOK || statusCode == http.StatusCreated {
		if ResponseBody != nil {
			err = json.Unmarshal(data, ResponseBody)
			if err != nil {
				return &Response{StatusCode: statusCode, Body: data},
					fmt.Errorf("Error parsing response body of %s request to [%s]: %s", Method, URL, err)
			}
		}
		return &Response{StatusCode: statusCode, Body: data}, nil
	} else {
		var response Response
		json.Unmarshal(data, &response)
		// We're ignoring the error because the server might have returnes
		// a generic error message. We will store that message in the RawBody
		// and just go ahead
		response.StatusCode = statusCode
		response.Body = data
		return &response, nil
	}
}
