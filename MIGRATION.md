# Migrate from CIS to Person API

`mozilla/cis` is deprecated. Its successor is
[`mozilla/personapi`](https://github.com/mozilla/terraform-provider-personapi),
available starting at version `2.0.0`. Existing CIS releases remain available.

## Update configuration and references

Update provider requirements in root and child modules, provider blocks, aliases,
explicit provider selections, and module provider mappings.

| Old name | New name |
| --- | --- |
| `mozilla/cis` | `mozilla/personapi` |
| `provider "cis"` | `provider "personapi"` |
| `cis_people` | `personapi_people` |
| `cis_group` | `personapi_group` |

For Terraform:

```hcl
terraform {
  required_providers {
    personapi = {
      source  = "registry.terraform.io/mozilla/personapi"
      version = "~> 2.0"
    }
  }
}

provider "personapi" {}

data "personapi_people" "alex" {
  email = "alex@mozilla.com"
}

data "personapi_group" "team" {
  ldap_group = "team_moco"
  staff      = true
}
```

For OpenTofu, use `registry.opentofu.org/mozilla/personapi` as the source.
Copy existing provider settings, aliases, and authentication configuration to the
corresponding successor provider configuration. Review the successor's documentation
for its supported settings.
