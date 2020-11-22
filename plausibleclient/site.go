package plausibleclient

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) CreateSite(domain, timezone string) (*SiteSettings, error) {
	if !c.loggedIn {
		err := c.login()
		if err != nil {
			return nil, err
		}
	}

	doc, err := c.getDocument("/sites/new")
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

	doc, err := c.getDocument("/" + domain + "/settings/danger-zone")
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
	_, err = c.httpClient.PostForm("https://plausible.io/"+domain, values)
	return err
}

func (c *Client) UpdateSite(domain, timezone string) (*SiteSettings, error) {
	if !c.loggedIn {
		err := c.login()
		if err != nil {
			return nil, err
		}
	}

	doc, err := c.getDocument("/" + domain + "/settings/general")
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
	Domain      string
	Timezone    string
	SharedLinks []string
	Goals       []int
}

func (c *Client) GetSiteSettings(domain string) (*SiteSettings, error) {
	if !c.loggedIn {
		err := c.login()
		if err != nil {
			return nil, err
		}
	}

	doc, err := c.getDocument("/" + domain + "/settings/general")
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

	doc, err = c.getDocument("/" + domain + "/settings/visibility")
	if err != nil {
		return nil, err
	}

	var sharedLinks []string
	doc.Find(`[value*='https://plausible.io/share/']`).Each(func(i int, s *goquery.Selection) {
		sharedLink, sharedLinkExists := s.Attr("value")
		if sharedLinkExists {
			sharedLinks = append(sharedLinks, sharedLink)
		}
	})

	doc, err = c.getDocument("/" + domain + "/settings/goals")
	if err != nil {
		return nil, err
	}

	var goals []int
	var errs []error
	doc.Find(`button[data-to*="/` + domain + `/goals/"]`).Each(func(i int, s *goquery.Selection) {
		g, exists := s.Attr("data-to")
		if exists {
			parts := strings.Split(g, "/")
			id, err := strconv.Atoi(parts[len(parts)-1])
			if err != nil {
				errs = append(errs, err)
				return
			}
			goals = append(goals, id)
		}
	})

	if len(errs) > 0 {
		return nil, fmt.Errorf("Could not parse goal ids: %v", errs)
	}

	return &SiteSettings{
		Domain:      domain,
		Timezone:    timezone,
		SharedLinks: sharedLinks,
		Goals:       goals,
	}, nil
}

func (c *Client) ListSites() ([]string, error) {
	if !c.loggedIn {
		err := c.login()
		if err != nil {
			return nil, err
		}
	}
	doc, err := c.getDocument("/sites")
	if err != nil {
		return nil, err
	}

	var domains []string
	doc.Find("a[href*='/settings']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		parts := strings.Split(href, "/")
		if exists {
			domains = append(domains, parts[1])
		}
	})

	return domains, nil
}
