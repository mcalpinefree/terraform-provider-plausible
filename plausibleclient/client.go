package plausibleclient

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Client struct {
	httpClient *http.Client
	username   string
	password   string
	loggedIn   bool
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

	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	c.httpClient = &http.Client{
		Jar: jar,
	}

	return &c
}
