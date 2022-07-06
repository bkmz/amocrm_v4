package amocrm_v4

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
)

type requestOpts struct {
	Method         string
	Path           string
	URLParameters  interface{}
	DataParameters interface{}
	Ret            interface{}
}

func httpRequest(opts requestOpts) error {
	var buf bytes.Buffer

	if opts.DataParameters != nil {
		if err := json.NewEncoder(&buf).Encode(&opts.DataParameters); err != nil {
			return err
		}
	}

	// set URL parameters
	values, err := query.Values(opts.URLParameters)
	if err != nil {
		return err
	}

	requestURL := client.getUrl(opts.Path)
	if len(values) > 0 {
		requestURL += "?" + values.Encode()
	}

	log.Debugf("Request URL: %s", requestURL)
	log.Debugf("URL Parameters: %s", values.Encode())
	log.Debugf("Body Parameters: %s", buf.String())

	req, err := http.NewRequest(opts.Method, client.getUrl(opts.Path), &buf)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	if opts.Path != "/oauth2/access_token" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.accessToken))
	}

	log.Debugf("Request Headers: %s", req.Header)
	log.Debugf("Request: %+v", req)

	resp, err := client.client.Do(req)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("Ошибка при закрытии потока ответа: %v", err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Debugf("Response Body: %s", string(body))

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	err = json.Unmarshal(body, &opts.Ret)
	if err != nil {
		return err
	}

	return nil
}
