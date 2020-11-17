---
page_title: "plausible_goal Resource - terraform-provider-plausible"
subcategory: ""
description: |-
  
---

# Resource `plausible_goal`





## Schema

### Required

- **site_id** (String, Required) The domain of the site to create the goal for.

### Optional

- **event_name** (String, Optional) Custom event E.g. `Signup`
- **page_path** (String, Optional) Page path event. E.g. `/success`

### Read-only

- **id** (String, Read-only) The goal ID


