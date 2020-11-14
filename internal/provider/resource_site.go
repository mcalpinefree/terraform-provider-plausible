package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mcalpinefree/terraform-provider-plausible/plausibleclient"
)

func resourceSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteCreate,
		Read:   resourceSiteRead,
		Update: resourceSiteUpdate,
		Delete: resourceSiteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
	siteSettings, err := client.plausibleClient.CreateSite(domain, timezone)
	if err != nil {
		return err
	}
	d.SetId(siteSettings.Domain)
	return resourceSiteSetResourceData(siteSettings, d)
}

func resourceSiteSetResourceData(siteSettings *plausibleclient.SiteSettings, d *schema.ResourceData) error {
	d.Set("domain", siteSettings.Domain)
	d.Set("timezone", siteSettings.Timezone)
	d.Set("javascript_snippet", fmt.Sprintf(`<script async defer data-domain="%s" src="https://plausible.io/js/plausible.js"></script>`, siteSettings.Domain))
	return nil
}

func resourceSiteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient)
	id := d.Id()

	siteSettings, err := client.plausibleClient.GetSiteSettings(id)
	if err != nil {
		return err
	}

	return resourceSiteSetResourceData(siteSettings, d)
}

func resourceSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	// use the meta value to retrieve your client from the provider configure method
	client := meta.(*apiClient)
	id := d.Id()

	timezone := d.Get("timezone").(string)
	siteSettings, err := client.plausibleClient.UpdateSite(id, timezone)
	if err != nil {
		return err
	}

	return resourceSiteSetResourceData(siteSettings, d)
}

func resourceSiteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient)
	domain := d.Id()
	return client.plausibleClient.DeleteSite(domain)
}
