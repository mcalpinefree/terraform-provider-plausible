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

- **password** (String, Optional) Add a password or leave it blank so anyone with the link can see the stats.

### Read-only

- **id** (String, Read-only) The shared link ID
- **link** (String, Read-only) Shared link


