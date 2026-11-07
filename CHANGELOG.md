# Changelog

## 1.7.0 (2026-11-07)

### Deprecated

- `mozilla/cis` has been superseded by [`mozilla/personapi`](https://github.com/mozilla/terraform-provider-personapi). This is the final planned CIS release; existing releases remain available.
- The provider, `cis_people`, and `cis_group` now emit deprecation warnings. Their configuration and lookup behavior are unchanged.
- See the [migration guide](MIGRATION.md) before switching providers.
