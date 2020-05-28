package notification

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func CreateComment(token, githubURL, issue, comment string) error {
	endpoint := fmt.Sprintf("https://api.%s/issues/%s/comments", githubURL, issue)

	header := map[string]string{
		"Authorization": token,
	}

	type Comment struct {
		Body string `json:"body"`
	}

	c := &Comment{
		Body: comment,
	}

	jsonValue, err := json.Marshal(c)
	if err != nil {
		return err
	}

	var res interface{}

	if err := postBody(endpoint, header, nil, jsonValue, res); err != nil {
		return err
	}
	return nil
}

// postBody makes a post request with json body
func postBody(endpoint string, header, query map[string]string, jsonValue []byte, res interface{}) error {
	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	return handleRequest(request, header, query, res)
}

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
		json, err := json.Marshal(resp.Body)
		if err != nil {
			log.Println(err)
		}
		return errors.New(string(json))
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
