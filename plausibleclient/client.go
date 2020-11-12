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

func (c *Client) CreateSite(domain, timezone string) (*SiteSettings, error) {
	if !c.loggedIn {
		err := c.login()
		if err != nil {
			return nil, err
		}
	}
	// get csrf token
	resp, err := c.httpClient.Get("https://plausible.io/sites/new")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Find the form CSRF token
	csrfToken := ""
	csrfTokenExists := false
	doc.Find(`form > input[name="_csrf_token"]`).Each(func(i int, s *goquery.Selection) {
		csrfToken, csrfTokenExists = s.Attr("value")
	})
	if !csrfTokenExists {
		return nil, fmt.Errorf("could not find csrf token in HTML form")
	}

	values := url.Values{}
	values.Add("_csrf_token", csrfToken)
	values.Add("site[domain]", domain)
	values.Add("site[timezone]", timezone)
	_, err = c.httpClient.PostForm("https://plausible.io/sites", values)
	return &SiteSettings{
		Domain:   domain,
		Timezone: timezone,
	}, err
}

func (c *Client) DeleteSite(domain string) error {
	if !c.loggedIn {
		err := c.login()
		if err != nil {
			return err
		}
	}
	// get csrf token
	resp, err := c.httpClient.Get("https://plausible.io/" + domain + "/settings")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	// Find the form CSRF token
	csrfToken := ""
	csrfTokenExists := false
	doc.Find(`a[data-method="delete"]`).Each(func(i int, s *goquery.Selection) {
		csrfToken, csrfTokenExists = s.Attr("data-csrf")
	})
	if !csrfTokenExists {
		return fmt.Errorf("could not find csrf token in HTML form")
	}
	values := url.Values{}
	values.Add("_csrf_token", csrfToken)
	values.Add("_method", "delete")
	resp, err = c.httpClient.PostForm("https://plausible.io/"+domain, values)
	return err
}

func (c *Client) UpdateSite(domain, timezone string) (*SiteSettings, error) {
	if !c.loggedIn {
		err := c.login()
		if err != nil {
			return nil, err
		}
	}
	// get csrf token
	resp, err := c.httpClient.Get("https://plausible.io/" + domain + "/settings")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Find the form CSRF token
	csrfToken := ""
	csrfTokenExists := false
	doc.Find(`form > input[name="_csrf_token"]`).Each(func(i int, s *goquery.Selection) {
		csrfToken, csrfTokenExists = s.Attr("value")
	})
	if !csrfTokenExists {
		return nil, fmt.Errorf("could not find csrf token in HTML form")
	}

	values := url.Values{}
	values.Add("_csrf_token", csrfToken)
	values.Add("_method", "put")
	values.Add("site[timezone]", timezone)
	_, err = c.httpClient.PostForm("https://plausible.io/"+domain+"/settings", values)
	return &SiteSettings{
		Domain:   domain,
		Timezone: timezone,
	}, err
}

type SiteSettings struct {
	Domain   string
	Timezone string
}

func (c *Client) GetSiteSettings(domain string) (*SiteSettings, error) {
	if !c.loggedIn {
		err := c.login()
		if err != nil {
			return nil, err
		}
	}
	resp, err := c.httpClient.Get("https://plausible.io/" + domain + "/settings")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	domainExists := false
	doc.Find("#site_domain").Each(func(i int, s *goquery.Selection) {
		domain, domainExists = s.Attr("value")
	})
	if !domainExists {
		return nil, fmt.Errorf("could not find domain in HTML document for %s", domain)
	}

	timezone := ""
	timezoneExists := false
	doc.Find(`#site_timezone > option[selected=""]`).Each(func(i int, s *goquery.Selection) {
		timezone, timezoneExists = s.Attr("value")
	})
	if !timezoneExists {
		return nil, fmt.Errorf("could not find timezone in HTML document for %s", domain)
	}

	return &SiteSettings{
		Domain:   domain,
		Timezone: timezone,
	}, nil
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
