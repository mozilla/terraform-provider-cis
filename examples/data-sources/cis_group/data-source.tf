terraform {
  required_providers {
    cis = {
      source = "registry.opentofu.org/mozilla/cis"
    }
  }
}

provider "cis" {}

data "cis_group" "example" {
  ldap_group = "team_moco"
  staff      = true
}

output "example" {
  value = data.cis_group.example
}
