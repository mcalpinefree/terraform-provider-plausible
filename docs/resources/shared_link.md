---
page_title: "plausible_shared_link Resource - terraform-provider-plausible"
subcategory: ""
description: |-
  
---

# Resource `plausible_shared_link`





## Schema

### Required

- **site_id** (String, Required) The domain of the site to create the shared link for.

### Optional

- **id** (String, Optional) The ID of this resource.
- **password** (String, Optional) Add a password or leave it blank so anyone with the link can see the stats.

### Read-only

- **link** (String, Read-only) Shared link


