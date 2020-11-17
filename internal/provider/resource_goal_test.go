package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceGoal(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGoal(acctest.RandomWithPrefix("testacc-tf")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("plausible_goal.testacc", "id", regexp.MustCompile(`^[0-9]+$`)),
				),
			},
		},
	})
}

func testAccResourceGoal(domain string) string {
	return fmt.Sprintf(`
resource "plausible_site" "testacc" {
    domain = "%s"
    timezone = "Pacific/Auckland"
}

resource "plausible_goal" "testacc" {
  site_id   = plausible_site.testacc.id
  page_path = "/success"
}
	`, domain)
}
