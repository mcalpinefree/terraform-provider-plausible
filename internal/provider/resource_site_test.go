package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSite(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSite(acctest.RandomWithPrefix("testacc-tf")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("plausible_site.testacc", "timezone", "Pacific/Auckland"),
				),
			},
		},
	})
}

func testAccResourceSite(domain string) string {
	return fmt.Sprintf(`
resource "plausible_site" "testacc" {
    domain = "%s"
    timezone = "Pacific/Auckland"
}
	`, domain)
}
