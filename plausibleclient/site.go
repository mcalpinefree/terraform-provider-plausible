package plausibleclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func (c *Client) GetSite(domain string) (*Site, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/v1/sites/"+domain, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	site := Site{}
	err = json.Unmarshal(b, &site)
	if err != nil {
		return nil, err
	}
	return &site, nil
}

func (c *Client) CreateSite(domain, timezone string) (*Site, error) {
	values := url.Values{}
	values.Add("domain", domain)
	values.Add("timezone", timezone)
	resp, err := c.postForm(c.baseURL+"/api/v1/sites", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	s := Site{}
	err = json.Unmarshal(b, &s)
	return &s, err
}

func (c *Client) DeleteSite(domain string) error {
	log.Printf("[DEBUG] DeleteSite")

	req, err := http.NewRequest("DELETE", c.baseURL+"/api/v1/sites/"+domain, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	a := struct {
		Deleted bool `json:"deleted"`
	}{}
	err = json.Unmarshal(b, &a)
	if err != nil {
		return err
	}

	if !a.Deleted {
		return fmt.Errorf("could not delete site %s", domain)
	}

	return nil
}
