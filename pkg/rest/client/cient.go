package client

import (
	"bytes"
	"encoding/json"
	"errors"
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

type errorMessage struct {
	Message string `json:"message"`
}

func Get(URL string, ResponseBody interface{}) (int, error) {
	return PerformRequest("GET", URL, nil, ResponseBody)
}

func Post(URL string, RequestBody interface{}, ResponseBody interface{}) (int, error) {
	return PerformRequest("POST", URL, RequestBody, ResponseBody)
}

func Put(URL string, RequestBody interface{}, ResponseBody interface{}) (int, error) {
	return PerformRequest("PUT", URL, RequestBody, ResponseBody)
}
func Delete(URL string) (int, error) {
	return PerformRequest("DELETE", URL, nil, nil)
}

func PerformRequest(Method string, URL string, RequestBody interface{}, ResponseBody interface{}) (int, error) {
	var statusCode int
	if Method == "" {
		return statusCode, fmt.Errorf("No method provided for HTTP request")
	}
	if URL == "" {
		return statusCode, fmt.Errorf("No URL provided for HTTP request")
	}
	var reader io.Reader
	if RequestBody != nil {
		buf := new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(RequestBody)
		if err != nil {
			return statusCode, fmt.Errorf("Error marshalling request body for %s request to [%s]: %s", Method, URL, err)
		}
		reader = buf
	}
	req, err := http.NewRequest(Method, URL, reader)
	if err != nil {
		return statusCode, fmt.Errorf("Error creating %s request equest to [%s]: %s", Method, URL, err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	res, err := httpClient.Do(req)
	if res != nil {
		statusCode = res.StatusCode
	}
	if err != nil {
		return statusCode, fmt.Errorf("Error performing %s request equest to [%s]: %s", Method, URL, err)
	}
	defer res.Body.Close()
	// We need to read all because else keep-alive won't work
	data, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(data))
	if err != nil {
		return statusCode, fmt.Errorf("Error reading response body of %s request to [%s]: %s", Method, URL, err)
	}
	if ResponseBody != nil && (statusCode == http.StatusOK || statusCode == http.StatusCreated) {
		err = json.Unmarshal(data, ResponseBody)
		if err != nil {
			return statusCode, fmt.Errorf("Error parsing response body of %s request to [%s]: %s", Method, URL, err)
		}
	}
	if statusCode > 400 {
		var errMsg errorMessage
		err = json.Unmarshal(data, &errMsg)
		if err != nil {
			return statusCode, errors.New(errMsg.Message)
		}
	}
	return statusCode, nil
}
