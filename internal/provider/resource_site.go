package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteCreate,
		Read:   resourceSiteRead,
		Update: resourceSiteUpdate,
		Delete: resourceSiteDelete,

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"javascript_snippet": {
				Description: "Include this snippet in the <head> of your website.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSiteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient)
	domain := d.Get("domain").(string)
	timezone := d.Get("timezone").(string)
	err := client.plausibleClient.CreateSite(domain, timezone)
	if err != nil {
		return err
	}
	d.SetId(domain)
	d.Set("javascript_snippet", fmt.Sprintf(`<script async defer data-domain="%s" src="https://plausible.io/js/plausible.js"></script>`, domain))
	return nil
}

func resourceSiteRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	timezone := d.Get("timezone").(string)

	d.Set("domain", id)
	d.Set("timezone", timezone)
	d.Set("javascript_snippet", fmt.Sprintf(`<script async defer data-domain="%s" src="https://plausible.io/js/plausible.js"></script>`, id))

	return nil
}

func resourceSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*apiClient)
	id := d.Id()

	timezone := d.Get("timezone").(string)
	err := client.plausibleClient.UpdateSite(id, timezone)
	if err != nil {
		return err
	}
	d.Set("timezone", timezone)

	return nil
}

func resourceSiteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient)
	domain := d.Get("domain").(string)
	return client.plausibleClient.DeleteSite(domain)
}
