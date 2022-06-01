package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mcalpinefree/terraform-provider-plausible/plausibleclient"
)

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"url": {
					Description: "Plausible URL. Can be specified with the `PLAUSIBLE_URL` environment variable.",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("PLAUSIBLE_URL", "https://plausible.io"),
				},
				"api_key": {
					Description: "Plausible API KEY. Can be specified with the `PLAUSIBLE_API_KEY` environment variable.",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("PLAUSIBLE_API_KEY", ""),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				//"scaffolding_data_source": dataSourceScaffolding(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"plausible_site": resourceSite(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

type apiClient struct {
	// Add whatever fields, client or connection info, etc. here
	// you would need to setup to communicate with the upstream
	// API.
	plausibleClient *plausibleclient.Client
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// Setup a User-Agent for your API client (replace the provider name for yours):
		// userAgent := p.UserAgent("terraform-provider-plausible", version)
		// TODO: myClient.UserAgent = userAgent
		url := d.Get("url").(string)
		apiKey := d.Get("api_key").(string)
		c := plausibleclient.NewClient(url, apiKey)
		return &apiClient{plausibleClient: c}, nil
	}
}
