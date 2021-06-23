package restclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func handleRequest(request *http.Request, header map[string]string, query map[string]string, res interface{}) error {
	for key, value := range header {
		request.Header.Add(key, value)
	}

	q := request.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}

	request.URL.RawQuery = q.Encode()

	client := &http.Client{Timeout: time.Duration(150 * time.Second)}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		return errors.New(string(b))
	}

	// b, _ := ioutil.ReadAll(resp.Body)
	// log.Println(string(b))
	// This was a delete request and was successful.
	// no need to try decode the body.
	if resp.StatusCode == 204 {
		return nil
	}

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	return nil
}

// FileUpload takes in data and handles making the put request
func FileUpload(endpoint string, header, query map[string]string, body io.Reader, res interface{}) error {
	request, err := http.NewRequest(http.MethodPut, endpoint, body)
	if err != nil {
		return err
	}

	return handleRequest(request, header, query, res)
}

// PostForm makes a post request with form data
func PostForm(endpoint string, header, query map[string]string, form url.Values, res interface{}) error {
	request, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	return handleRequest(request, header, query, res)
}

// PostBody makes a post request with json body
func PostBody(endpoint string, header, query map[string]string, jsonValue []byte, res interface{}) error {
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	return handleRequest(request, header, query, res)
}

// Put makes a put request
func Put(endpoint string, header, query map[string]string, body []byte, res interface{}) error {
	request, err := http.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	return handleRequest(request, header, query, res)
}

// Patch makes a patch request to url
func Patch(endpoint string, header, query map[string]string, body []byte, res interface{}) error {
	request, err := http.NewRequest(http.MethodPatch, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return handleRequest(request, header, query, res)
}

// Fetch makes a get request to endpoint
func Fetch(endpoint string, header, query map[string]string, res interface{}) error {
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)

	if err != nil {
		return err
	}

	return handleRequest(request, header, query, res)
}

func Delete(endpoint string, header, query map[string]string, res interface{}) error {
	request, err := http.NewRequest(http.MethodDelete, endpoint, nil)

	if err != nil {
		return err
	}

	return handleRequest(request, header, query, res)
}

func Head(endpoint string, header, query map[string]string, res interface{}) error {
	request, err := http.NewRequest(http.MethodHead, endpoint, nil)

	if err != nil {
		return err
	}

	return handleRequest(request, header, query, res)
}
