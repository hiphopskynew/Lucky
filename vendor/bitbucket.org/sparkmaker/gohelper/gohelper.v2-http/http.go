package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"bitbucket.org/sparkmaker/gohelper/logger/stdout"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PUTCH  = "PUTCH"
	PATCH  = "PATCH"
	DELETE = "DELETE"
)

var (
	defaultTimeOut   = 30 * time.Second
	defaultMaxRetry  = 0
	defaultTimeDelay = 5 * time.Second
)

type Property struct {
	URL         string
	Body        interface{}
	Headers     map[string]string
	AttempRetry *AttempRetry
	PrintLog    bool
}

type HttpCaller interface {
	GET() (string, error)
	POST() (string, error)
	PUT() (string, error)
	PATCH() (string, error)
	DELETE() (string, error)
}

type caller struct {
	URL         string
	Method      string
	Body        interface{}
	Headers     map[string]string
	AttempRetry *AttempRetry
	Printlog    bool
}

func retry(client *http.Client, req *http.Request, c *caller) (*http.Response, error) {
	for c.AttempRetry.attempRetry <= c.AttempRetry.MaxRetry {
		c.printLogInfo(fmt.Sprintln("call round", c.AttempRetry.attempRetry))
		var isRetry bool
		resp, err := client.Do(req)
		if err != nil {
			c.printLogError(err)
			isRetry = true
		}
		if resp != nil {
			// run
			c.AttempRetry.Response = resp
			for _, retryFn := range c.AttempRetry.RetryFns {
				isRetry = retryFn(*c.AttempRetry)
				if isRetry {
					break
				}
			}
		}

		if !isRetry || c.AttempRetry.attempRetry == c.AttempRetry.MaxRetry {
			c.printLogInfo("out of retry")
			if err != nil {
				c.printLogError(fmt.Sprintf("Server error message: %s", err))
				return nil, errors.New("http caller timeout")
			}
			return resp, err
		}
		select {
		case <-time.After(c.AttempRetry.TimeDelay):
			c.AttempRetry.attempRetry++
		}
	}

	return nil, errors.New("panic retry function")
}

func (c *caller) printLogInfo(msg interface{}) {
	if c.Printlog {
		stdout.Info(msg)
	}
}

func (c *caller) printLogError(msg interface{}) {
	if c.Printlog {
		stdout.Error(msg)
	}
}

func invoke(c *caller) (string, error) {
	body := bytes.NewBuffer(nil)
	if c.Body != nil {
		switch c.Body.(type) {
		case *bytes.Buffer:
			body = c.Body.(*bytes.Buffer)
		default:
			b, _ := json.Marshal(c.Body)
			body = bytes.NewBuffer(b)
		}
	}
	c.printLogInfo(fmt.Sprintf("Retry settings: %#v", c.AttempRetry))
	req, _ := http.NewRequest(c.Method, c.URL, body)
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{
		Timeout: c.AttempRetry.Timeout,
	}
	resp, err := retry(client, req, c)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, _ := ioutil.ReadAll(resp.Body)
	return string(content), nil
}

func (c *caller) GET() (string, error) {
	c.Method = GET
	return invoke(c)
}

func (c *caller) POST() (string, error) {
	c.Method = POST
	return invoke(c)
}

func (c *caller) PUT() (string, error) {
	c.Method = PUT
	return invoke(c)
}

func (c *caller) PATCH() (string, error) {
	c.Method = PATCH
	return invoke(c)
}

func (c *caller) DELETE() (string, error) {
	c.Method = DELETE
	return invoke(c)
}

func New(property Property) HttpCaller {
	//TODO set default for attemp retry
	c := &caller{
		URL:     property.URL,
		Body:    property.Body,
		Headers: property.Headers,
	}
	if property.AttempRetry != nil {
		c.AttempRetry = property.AttempRetry
	} else {
		c.AttempRetry = &AttempRetry{
			TimeDelay: defaultTimeDelay,
			Timeout:   defaultTimeOut,
			MaxRetry:  defaultMaxRetry,
		}
	}
	c.Printlog = property.PrintLog
	c.AttempRetry.attempRetry = 0

	return c
}
