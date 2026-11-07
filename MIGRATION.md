# Migrate from CIS to Person API

`mozilla/cis` is deprecated. Its successor is
[`mozilla/personapi`](https://github.com/mozilla/terraform-provider-personapi),
available starting at version `2.0.0`. Existing CIS releases remain available.

The successor was initially discussed as `mozillaiam`; its published name is `personapi`.

## 1. Back up state and remove old cached lookups

Perform the migration for each affected workspace. Before changing provider
configuration, initialize the existing configuration if necessary, inspect its state,
and save a backup in a secure location:

```sh
terraform state list
terraform state pull > before-personapi.tfstate
```

For OpenTofu, use `tofu` instead of `terraform` throughout this guide.
State backups contain identity information; keep them out of version control.

This provider exposes only data sources. Remove the old CIS data source entries
from state so the successor can read them again. Use the actual addresses from
`state list`, including module paths and instance keys. For example:

```sh
terraform state rm 'data.cis_people.alex' 'data.cis_group.team'
```

This forgets cached lookup results; it does not delete people, groups, or memberships.
Do not remove resources belonging to other providers.

## 2. Update configuration and references

Update provider requirements in root and child modules, provider blocks, aliases,
explicit provider selections, and module provider mappings.

| Old name | New name |
| --- | --- |
| `mozilla/cis` | `mozilla/personapi` |
| `provider "cis"` | `provider "personapi"` |
| `cis_people` | `personapi_people` |
| `cis_group` | `personapi_group` |
| `data.cis_people.alex` | `data.personapi_people.alex` |
| `data.cis_group.team` | `data.personapi_group.team` |

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

Changing a provider's local name or using `state replace-provider` alone does not
rename data source types. Removing the old cached data source entries in step 1
avoids retaining entries whose type names the successor does not support.

## 3. Initialize and review the plan

```sh
terraform init
terraform validate
terraform plan
```

Confirm that the successor is installed and the old CIS provider is no longer
required by configuration or state. Review all downstream changes: refreshed
identity and membership results can affect resources managed by other providers.
Apply the reviewed plan and commit configuration and `.terraform.lock.hcl` changes.

Repeat for every affected workspace. Keep the old provider's release assets intact;
there is no automatic registry redirect from `mozilla/cis` to `mozilla/personapi`.
