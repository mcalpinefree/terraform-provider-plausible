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
				"username": {
					Description: "Plausible username. Can be specified with the `PLAUSIBLE_USERNAME` environment variable.",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("PLAUSIBLE_USERNAME", ""),
				},
				"password": {
					Description: "Plausible password. Can be specified with the `PLAUSIBLE_PASSWORD` environment variable.",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("PLAUSIBLE_PASSWORD", ""),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"scaffolding_data_source": dataSourceScaffolding(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"scaffolding_resource": resourceScaffolding(),
				"plausible_site":       resourceSite(),
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
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		c := plausibleclient.NewClient(username, password)
		return &apiClient{plausibleClient: c}, nil
	}
}
