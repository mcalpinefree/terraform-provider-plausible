package plausibleclient

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Goal struct {
	ID        int
	Domain    string
	PagePath  *string
	EventName *string
}

type GoalType int

const (
	PagePath GoalType = iota
	EventName
)

func (c *Client) CreateGoal(domain string, goalType GoalType, goal string) (*Goal, error) {
	if !c.loggedIn {
		err := c.login()
		if err != nil {
			return nil, err
		}
	}

	c.mutexkv.Lock(domain)
	defer c.mutexkv.Unlock(domain)

	resp, err := c.httpClient.Get("https://plausible.io/" + domain + "/goals/new")
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
	result := Goal{Domain: domain}
	switch goalType {
	case PagePath:
		result.PagePath = &goal
	case EventName:
		result.EventName = &goal
	}

	if result.PagePath != nil {
		values.Add("goal[page_path]", *result.PagePath)
	} else {
		values.Add("goal[page_path]", "")
	}
	if result.EventName != nil {
		values.Add("goal[event_name]", *result.EventName)
	} else {
		values.Add("goal[event_name]", "")
	}

	before, err := c.GetSiteSettings(domain)
	if err != nil {
		return nil, err
	}

	_, err = c.httpClient.PostForm("https://plausible.io/"+domain+"/goals", values)
	if err != nil {
		return nil, err
	}

	after, err := c.GetSiteSettings(domain)
	if err != nil {
		return nil, err
	}

	if len(before.Goals) != (len(after.Goals) - 1) {
		return nil, fmt.Errorf("expected there to be one more goal after requesting to create a new one, but the count went from %d to %d", len(before.SharedLinks), len(after.SharedLinks))
	}

AFTER:
	for _, v := range after.Goals {
		for _, w := range before.Goals {
			if v == w {
				continue AFTER
			}
		}
		result.ID = v
		return &result, nil
	}

	return nil, fmt.Errorf("could not find newly created goal")
}

func (c *Client) GetGoal(domain string, id int) (*Goal, error) {
	if !c.loggedIn {
		err := c.login()
		if err != nil {
			return nil, err
		}
	}

	c.mutexkv.Lock(domain)
	defer c.mutexkv.Unlock(domain)

	result := Goal{ID: id, Domain: domain}

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

	cssSelector := fmt.Sprintf(`button[data-to="/%s/goals/%d"]`, domain, id)
	doc.Find(cssSelector).Each(func(i int, s *goquery.Selection) {
		var c string
		s.SiblingsFiltered("small").Each(func(i int, s *goquery.Selection) {
			h, _ := s.Html()
			c = strings.TrimSpace(h)
		})
		if strings.HasPrefix(c, "Visit /") {
			pagePath := strings.TrimPrefix(c, "Visit ")
			result.PagePath = &pagePath
		} else {
			result.EventName = &c
		}
	})

	return &result, nil
}

func (c *Client) DeleteGoal(domain string, id int) error {
	if !c.loggedIn {
		err := c.login()
		if err != nil {
			return err
		}
	}

	c.mutexkv.Lock(domain)
	defer c.mutexkv.Unlock(domain)

	resp, err := c.httpClient.Get("https://plausible.io/" + domain + "/settings")
	if err != nil {
		return err
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
	cssSelector := fmt.Sprintf(`button[data-to="/%s/goals/%d"]`, domain, id)
	doc.Find(cssSelector).Each(func(i int, s *goquery.Selection) {
		csrfToken, csrfTokenExists = s.Attr("data-csrf")
	})
	if !csrfTokenExists {
		return fmt.Errorf("could not find csrf token in HTML form")
	}

	values := url.Values{}
	values.Add("_csrf_token", csrfToken)
	values.Add("_method", "delete")
	_, err = c.httpClient.PostForm(fmt.Sprintf("https://plausible.io/%s/goals/%d", domain, id), values)
	return err
}
