package plausibleclient

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	httpClient *http.Client
	username   string
	password   string
	loggedIn   bool
	mutexkv    *MutexKV
	baseURL    string
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
	resp, err := c.httpClient.PostForm("https://plausible.io/login", values)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not login, received status: %s", resp.Status)
	}

	c.loggedIn = true
	return nil
}

func NewClient(username, password string) *Client {
	c := Client{}
	c.username = username
	c.password = password
	c.baseURL = "https://plausible.io"

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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
