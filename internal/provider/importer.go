package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func importDomainAndID(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if id := d.Id(); strings.Contains(id, ":") {
		importParts := strings.SplitN(id, ":", 2)
		d.SetId(importParts[1])
		d.Set("site_id", importParts[0])
	}
	return []*schema.ResourceData{d}, nil
}
