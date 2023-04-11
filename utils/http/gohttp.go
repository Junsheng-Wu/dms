package http

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/wxnacy/wgo/arrays"
)

func Req(URL, method string, tlsVerify bool, data *[]byte, header *map[string]string, timeout int, resp *[]byte) error {
	_timeout := time.Duration(timeout) * time.Second
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !tlsVerify,
		},
	}
	client := http.Client{
		Timeout:   _timeout,
		Transport: tr,
	}

	var req *http.Request
	var err error
	if arrays.ContainsString([]string{"GET", "DELETE"}, method) != -1 {
		req, err = http.NewRequest(method, URL, nil)
		if err != nil {
			return errors.Wrap(err, "http.NewRequest")
		}
	} else if arrays.ContainsString([]string{"PUT", "POST"}, method) != -1 {
		body := bytes.NewBuffer(*data)
		req, err = http.NewRequest(method, URL, body)
		if err != nil {
			return errors.Wrap(err, "http.NewRequest")
		}
	} else {
		return errors.New("Invalid method: " + method)
	}

	if header != nil {
		for key, value := range *header {
			req.Header.Set(key, value)
		}
	}

	res, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do")
	}
	defer res.Body.Close()

	*resp, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadAll")
	}

	if res.StatusCode >= http.StatusBadRequest {
		return errors.New(string(*resp))
	}

	return nil
}
