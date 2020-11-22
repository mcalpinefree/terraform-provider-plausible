package plausibleclient

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SharedLink struct {
	ID       string
	Domain   string
	Password string
	Link     string
}

func (c *Client) CreateSharedLink(domain, password string) (*SharedLink, error) {
	if !c.loggedIn {
		err := c.login()
		if err != nil {
			return nil, err
		}
	}

	c.mutexkv.Lock(domain)
	defer c.mutexkv.Unlock(domain)

	doc, err := c.getDocument("/sites/" + domain + "/shared-links/new")
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

	// super hacky but the /shared-links POST does not return the ID of the
	// shared link that was created so we have to look up the shared links
	// before and after the request and compare to find out what was
	// created
	before, err := c.GetSiteSettings(domain)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("_csrf_token", csrfToken)
	values.Add("shared_link[password]", password)
	_, err = c.httpClient.PostForm("https://plausible.io/sites/"+domain+"/shared-links", values)
	if err != nil {
		return nil, err
	}

	after, err := c.GetSiteSettings(domain)
	if err != nil {
		return nil, err
	}

	if len(before.SharedLinks) != (len(after.SharedLinks) - 1) {
		return nil, fmt.Errorf("expected there to be one more shared link after requesting to create a new one, but the count went from %d to %d", len(before.SharedLinks), len(after.SharedLinks))
	}

AFTER:
	for _, v := range after.SharedLinks {
		for _, w := range before.SharedLinks {
			if v == w {
				continue AFTER
			}
		}
		parts := strings.Split(v, "/")
		return &SharedLink{
			ID:       parts[len(parts)-1],
			Domain:   domain,
			Password: password,
			Link:     v,
		}, nil
	}

	return nil, fmt.Errorf("could not find newly created shared link")
}

func (c *Client) DeleteSharedLink(domain, id string) error {
	if !c.loggedIn {
		err := c.login()
		if err != nil {
			return err
		}
	}

	c.mutexkv.Lock(domain)
	defer c.mutexkv.Unlock(domain)

	doc, err := c.getDocument("/" + domain + "/settings/visibility")
	if err != nil {
		return err
	}

	// Find the form CSRF token
	csrfToken := ""
	csrfTokenExists := false
	cssSelector := fmt.Sprintf(`button[data-to="/sites/%s/shared-links/%s"]`, domain, id)
	log.Printf("[TRACE] cssSelector: %s\n", cssSelector)
	doc.Find(cssSelector).Each(func(i int, s *goquery.Selection) {
		csrfToken, csrfTokenExists = s.Attr("data-csrf")
	})
	if !csrfTokenExists {
		return fmt.Errorf("could not find csrf token in HTML form")
	}

	values := url.Values{}
	values.Add("_csrf_token", csrfToken)
	values.Add("_method", "delete")
	_, err = c.httpClient.PostForm("https://plausible.io/sites/"+domain+"/shared-links/"+id, values)
	return err
}
