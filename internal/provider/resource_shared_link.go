package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mcalpinefree/terraform-provider-plausible/plausibleclient"
)

func resourceSharedLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceSharedLinkCreate,
		Read:   resourceSharedLinkRead,
		Delete: resourceSharedLinkDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The shared link ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"site_id": {
				Description: "The domain of the site to create the shared link for.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the shared link to create.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"password": {
				Description: "Add a password or leave it blank so anyone with the link can see the stats.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"link": {
				Description: "Shared link",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSharedLinkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient)
	domain := d.Get("site_id").(string)
	name := d.Get("name").(string)
	password := d.Get("password").(string)
	sharedLink, err := client.plausibleClient.CreateSharedLink(domain, name, password)
	if err != nil {
		return err
	}
	d.SetId(sharedLink.ID)
	return resourceSharedLinkSetResourceData(sharedLink, d)
}

func resourceSharedLinkSetResourceData(sharedLink *plausibleclient.SharedLink, d *schema.ResourceData) error {
	d.Set("site_id", sharedLink.Domain)
	d.Set("password", sharedLink.Password)
	d.Set("link", sharedLink.Link)
	return nil
}

func resourceSharedLinkRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	sharedLink := &plausibleclient.SharedLink{
		ID:       id,
		Domain:   d.Get("site_id").(string),
		Password: d.Get("password").(string),
		Link:     d.Get("link").(string),
	}

	return resourceSharedLinkSetResourceData(sharedLink, d)
}

func resourceSharedLinkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient)
	id := d.Id()
	domain := d.Get("site_id").(string)
	return client.plausibleClient.DeleteSharedLink(domain, id)
}
