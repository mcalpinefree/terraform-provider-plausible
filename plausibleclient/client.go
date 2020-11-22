package plausibleclient

import (
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
	values := url.Values{}
	values.Add("email", c.username)
	values.Add("password", c.password)
	_, err := c.httpClient.PostForm("https://plausible.io/login", values)
	if err != nil {
		return err
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
	resp, err := c.httpClient.Get(c.baseURL + path)
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
