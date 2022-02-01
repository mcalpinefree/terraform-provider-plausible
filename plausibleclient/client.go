package plausibleclient

import (
	"fmt"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	httpClient  *http.Client
	username    string
	password    string
	loggedIn    bool
	mutexkv     *MutexKV
	baseURL     string
	maxAttempts int
}

func (c *Client) login() error {
	doc, err := c.getDocument("/login")
	if err != nil {
		return err
	}

	csrfToken := ""
	csrfTokenExists := false
	doc.Find(`form > input[name="_csrf_token"]`).Each(func(i int, s *goquery.Selection) {
		csrfToken, csrfTokenExists = s.Attr("value")
	})
	if !csrfTokenExists {
		return fmt.Errorf("could not find csrf token in login page")
	}

	values := url.Values{}
	values.Add("_csrf_token", csrfToken)
	values.Add("email", c.username)
	values.Add("password", c.password)
	resp, err := c.postForm(c.baseURL+"/login", values)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not login, received status: %s", resp.Status)
	}

	c.loggedIn = true
	return nil
}

func NewClient(url, username, password string) *Client {
	c := Client{}
	c.username = username
	c.password = password
	c.baseURL = url
	c.maxAttempts = 10

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	c.httpClient = &http.Client{
		Jar: jar,
	}

	c.mutexkv = NewMutexKV()

	return &c
}

func (c *Client) getDocument(path string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	attempts := 1
	for {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			if attempts > c.maxAttempts {
				return nil, fmt.Errorf("request failed after %d attempts with status: %s", c.maxAttempts, resp.Status)
			}
			time.Sleep(time.Duration(math.Pow(1.5, float64(attempts))) * time.Second)
			attempts++
			continue
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return nil, err
		}

		return doc, nil
	}
}

func (c *Client) postForm(url string, values url.Values) (*http.Response, error) {
	attempts := 1
	for {
		resp, err := c.httpClient.PostForm(url, values)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
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
