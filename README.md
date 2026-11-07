# CIS Terraform Provider (deprecated)

**Deprecated:** `mozilla/cis` has been superseded by
[`mozilla/personapi`](https://github.com/mozilla/terraform-provider-personapi).
Existing releases remain available. See the [migration guide](MIGRATION.md) before switching.

This provider reads people and group membership from Mozilla's IAM Person API.
The final planned CIS release adds deprecation warnings and migration documentation.
Existing provider configuration, `cis_people`, and `cis_group` retain their behavior.
New development continues in `mozilla/personapi`.

## Existing installations

Existing configurations can continue to use `mozilla/cis`. The provider documentation
is available in the [Terraform Registry](https://registry.terraform.io/providers/mozilla/cis/latest/docs)
and [OpenTofu Registry](https://search.opentofu.org/provider/mozilla/cis/latest).
Users pinned to releases before the deprecation release will not receive its runtime warnings.

## Development

Use the Go version specified in `go.mod`. Run `make test` for unit tests and
`make generate` to regenerate documentation. Documentation notices are defined in the
provider and data source schema descriptions so regeneration preserves them.

Acceptance tests query the live API. Set `AUTH0_CLIENT_ID` and `AUTH0_CLIENT_SECRET`
for an authorized client, then run `make testacc`.

## Final release

Before tagging the final release, run the checks and acceptance tests for the commit
being released. Pushing a `v*` tag triggers the signed GitHub release workflow.
GoReleaser includes the deprecation notice and migration guide in the release notes.
Preserve all historical release assets and checksums so existing installations can
continue to download and verify them.
