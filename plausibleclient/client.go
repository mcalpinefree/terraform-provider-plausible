package plausibleclient

import (
	"fmt"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	httpClient  *http.Client
	apiKey      string
	baseURL     string
	maxAttempts int
}

func NewClient(url, apiKey string) *Client {
	c := Client{}
	c.apiKey = apiKey
	c.baseURL = url
	c.maxAttempts = 10

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	c.httpClient = &http.Client{
		Jar: jar,
	}

	return &c
}

func (c *Client) postForm(url string, values url.Values) (*http.Response, error) {
	attempts := 1
	for {
		req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Authorization", "Bearer "+c.apiKey)
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			if attempts > c.maxAttempts {
				return nil, fmt.Errorf("request failed after %d attempts with status: %s", c.maxAttempts, resp.Status)
			}
			time.Sleep(time.Duration(math.Pow(1.5, float64(attempts))) * time.Second)
			attempts++
			continue
		}
		return resp, err

	}
}
