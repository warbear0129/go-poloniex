package poloniex

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type client struct {
	apiKey     string
	apiSecret  string
	httpClient *http.Client
	throttle   <-chan time.Time
}

const (
	DEFAULT_HTTPCLIENT_TIMEOUT = 30
)

var (
	reqInterval = 170 * time.Millisecond
)

// NewClient return a new Poloniex HTTP client
func NewClient(apiKey, apiSecret string) (c *client) {
	return &client{apiKey, apiSecret, &http.Client{}, time.Tick(reqInterval)}
}

// doTimeoutRequest do a HTTP request with timeout
func (c *client) doTimeoutRequest(req *http.Request) (*http.Response, error) {
	timeout := time.NewTimer(DEFAULT_HTTPCLIENT_TIMEOUT * time.Second)

	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		resp, err := c.httpClient.Do(req)
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timeout.C:
		return nil, errors.New("timeout on reading data from Poloniex API")
	}
}

func (c *client) makeRequest(method, command string, args map[string]string, respCh chan<- []byte, errCh chan<- error) {
	var req *http.Request
	body := []byte{}

	data := url.Values{}
	data.Add("command", command)
	if args != nil {
		for k, v := range args {
			data.Add(k, v)
		}
	}

	if method == "GET" {

		payload := data.Encode()
		reqURL := "https://poloniex.com/public?" + payload
		req, _ = http.NewRequest(method, reqURL, nil)

		fmt.Println(reqURL)

	} else if method == "POST" {

		if len(c.apiKey) == 0 || len(c.apiSecret) == 0 {
			respCh <- nil
			errCh <- errors.New("You need to set API Key and API Secret to call this method")
			return
		}

		reqURL := "https://poloniex.com/tradingApi"
		data.Add("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
		payload := data.Encode()

		mac := hmac.New(sha512.New, []byte(c.apiSecret))
		mac.Write([]byte(payload))
		sign := hex.EncodeToString(mac.Sum(nil))

		req, _ := http.NewRequest("POST", reqURL, strings.NewReader(payload))
		req.Header.Add("Key", c.apiKey)
		req.Header.Add("Sign", sign)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	}

	resp, err := c.doTimeoutRequest(req)
	if err != nil {
		respCh <- body
		errCh <- err
		return
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		respCh <- body
		errCh <- err
		return
	}
	if resp.StatusCode != 200 {
		respCh <- body
		errCh <- errors.New(resp.Status)
		return
	}

	respCh <- body
	errCh <- nil
	close(respCh)
	close(errCh)
}

// do prepare and process HTTP request to Poloniex API
func (c *client) do(method, command string, args map[string]string) (response []byte, err error) {
	respCh := make(chan []byte)
	errCh := make(chan error)
	<-c.throttle
	go c.makeRequest(method, command, args, respCh, errCh)
	response = <-respCh
	err = <-errCh
	return
}
