package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpCaller interface {
	GET() (string, error)
	POST() (string, error)
	PUT() (string, error)
	PATCH() (string, error)
	DELETE() (string, error)
	SetAttributes(URL string, Body interface{}, Headers map[string]string, RetryCount int, RetryDelay time.Duration) HttpCaller
}

type Caller struct {
	URL        string
	Body       interface{}
	Headers    map[string]string
	RetryCount int
	RetryDelay time.Duration
}

type value struct {
	URL        string
	Method     string
	Body       interface{}
	Headers    map[string]string
	RetryCount int
	RetryDelay time.Duration
}

func (caller Caller) SetAttributes(URL string, Body interface{}, Headers map[string]string, RetryCount int, RetryDelay time.Duration) HttpCaller {
	return Caller{
		URL:        URL,
		Body:       Body,
		Headers:    Headers,
		RetryCount: RetryCount,
		RetryDelay: RetryDelay,
	}
}

func (caller Caller) GET() (string, error) {
	return invoke(value{
		URL:        caller.URL,
		Method:     "GET",
		Headers:    caller.Headers,
		RetryCount: caller.RetryCount,
		RetryDelay: caller.RetryDelay,
	})
}

func (caller Caller) POST() (string, error) {
	return invoke(value{
		URL:        caller.URL,
		Method:     "POST",
		Body:       caller.Body,
		Headers:    caller.Headers,
		RetryCount: caller.RetryCount,
		RetryDelay: caller.RetryDelay,
	})
}

func (caller Caller) PUT() (string, error) {
	return invoke(value{
		URL:        caller.URL,
		Method:     "PUT",
		Body:       caller.Body,
		Headers:    caller.Headers,
		RetryCount: caller.RetryCount,
		RetryDelay: caller.RetryDelay,
	})
}

func (caller Caller) PATCH() (string, error) {
	return invoke(value{
		URL:        caller.URL,
		Method:     "PATCH",
		Body:       caller.Body,
		Headers:    caller.Headers,
		RetryCount: caller.RetryCount,
		RetryDelay: caller.RetryDelay,
	})
}

func (caller Caller) DELETE() (string, error) {
	return invoke(value{
		URL:        caller.URL,
		Method:     "DELETE",
		Body:       caller.Body,
		Headers:    caller.Headers,
		RetryCount: caller.RetryCount,
		RetryDelay: caller.RetryDelay,
	})
}

func retry(client *http.Client, req *http.Request, count int, delay time.Duration) (*http.Response, error) {
	// Fixed timeout 25 seconds can will configuration in the future
	timeout := 25

	if count < 0 {
		count = 0
	}

	if delay < time.Second {
		delay = time.Second
	}

	client.Timeout = time.Duration(timeout/(count+1)) * time.Second
	for i := 0; i <= count; i++ {
		resp, err := client.Do(req)
		if err != nil {
			// ignore delay
			// time.Sleep(delay)
			continue
		} else {
			return resp, err
		}
	}

	return nil, errors.New("HTTP caller timeout.")
}

func invoke(cp value) (string, error) {
	body := bytes.NewBuffer(nil)
	if cp.Body != nil {
		b, _ := json.Marshal(cp.Body)
		body = bytes.NewBuffer(b)
	}

	req, _ := http.NewRequest(cp.Method, cp.URL, body)
	for k, v := range cp.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := retry(client, req, cp.RetryCount, cp.RetryDelay)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	return string(content), nil
}
