package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSharedLink(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSharedLink(acctest.RandomWithPrefix("testacc-tf")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("plausible_shared_link.testacc", "link", regexp.MustCompile(`^https:\/\/plausible.io\/share\/[a-zA-Z0-9-_]+\?auth=[a-zA-Z0-9-_]+$`)),
				),
			},
		},
	})
}

func testAccResourceSharedLink(domain string) string {
	return fmt.Sprintf(`
resource "plausible_site" "testacc" {
    domain = "%s"
    timezone = "Pacific/Auckland"
}

resource "plausible_shared_link" "testacc" {
    site_id = plausible_site.testacc.id
}
	`, domain)
}
