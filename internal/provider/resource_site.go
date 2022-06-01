package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mcalpinefree/terraform-provider-plausible/plausibleclient"
)

func resourceSite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSiteCreate,
		ReadContext:   resourceSiteRead,
		DeleteContext: resourceSiteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The site ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"javascript_snippet": {
				Description: "Include this snippet in the <head> of your website.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSiteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	domain := d.Get("domain").(string)
	timezone := d.Get("timezone").(string)
	site, err := client.plausibleClient.CreateSite(domain, timezone)
	if err != nil {
		return diag.Errorf("error creating site (%s): %s", d.Id(), err)
	}
	d.SetId(site.Domain)
	return resourceSiteSetResourceData(site, d)
}

func resourceSiteSetResourceData(s *plausibleclient.Site, d *schema.ResourceData) diag.Diagnostics {
	d.Set("domain", s.Domain)
	d.Set("timezone", s.Timezone)
	d.Set("javascript_snippet", fmt.Sprintf(`<script defer data-domain="%s" src="https://plausible.io/js/plausible.js"></script>`, s.Domain))
	return nil
}

func resourceSiteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	site, err := client.plausibleClient.GetSite(d.Id())
	if err != nil {
		return diag.Errorf("error getting site (%s): %s", d.Id(), err)
	}
	return resourceSiteSetResourceData(site, d)
}

func resourceSiteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	domain := d.Id()
	err := client.plausibleClient.DeleteSite(domain)
	if err != nil {
		return diag.Errorf("error deleting site (%s): %s", d.Id(), err)
	}
	return nil
}
