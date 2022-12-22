package everypay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

var (
	client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:    10 * time.Second,
			MaxIdleConns:           0,
			MaxIdleConnsPerHost:    4,
			MaxConnsPerHost:        0,
			IdleConnTimeout:        90 * time.Second,
			ResponseHeaderTimeout:  2 * time.Second,
			ExpectContinueTimeout:  1 * time.Second,
			MaxResponseHeaderBytes: 2 * 1024,
			ForceAttemptHTTP2:      false,
		},
	}
)

type Everypay struct {
	endpoint *url.URL
	username string
	secret   string
	account  string
}

func NewEverypay(username, secret, account string, production bool) *Everypay {
	e := &Everypay{
		username: username,
		secret:   secret,
		account:  account,
	}

	if production {
		// https://pay.every-pay.eu/api/v4/
		e.endpoint = &url.URL{
			Scheme: "https",
			Host:   "pay.every-pay.eu",
			Path:   "/api/v4/",
		}
	} else {
		// https://igw-demo.every-pay.com/api/v4/
		e.endpoint = &url.URL{
			Scheme: "https",
			Host:   "igw-demo.every-pay.com",
			Path:   "/api/v4/",
		}
	}

	return e
}

func (e *Everypay) request(
	method string,
	reference *url.URL,
	requestData interface{},
	responseData interface{},
) error {
	var request *http.Request

	switch method {
	case http.MethodGet:
		request, _ = http.NewRequest(http.MethodGet, "", nil)

	case http.MethodPost:
		requestBody := &bytes.Buffer{}
		_ = json.NewEncoder(requestBody).Encode(requestData)

		request, _ = http.NewRequest(http.MethodPost, "", requestBody)
		request.Header.Set("Content-Type", "application/json; charset=utf-8")

	default:
		return fmt.Errorf("invalid method: %s", method)
	}

	request.URL = e.endpoint.ResolveReference(reference)
	request.SetBasicAuth(e.username, e.secret)

	response, err := client.Do(request)
	if err != nil {
		if response != nil {
			response.Body.Close()
		}
		return fmt.Errorf("send: %w", err)
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		break
	default:
		return fmt.Errorf("response: %s", response.Status)
	}

	if err := json.NewDecoder(response.Body).Decode(responseData); err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}
