package provider

const migrationGuideURL = "https://github.com/mozilla/terraform-provider-cis/blob/main/MIGRATION.md"

const deprecationMessage = "mozilla/cis is deprecated and has been superseded by mozilla/personapi. " +
	"Existing releases remain available. Migrate using the guide at " + migrationGuideURL

const deprecationNotice = "**Deprecated:** This provider has been superseded by " +
	"[`mozilla/personapi`](https://github.com/mozilla/terraform-provider-personapi). " +
	"Existing releases remain available. See the [migration guide](" + migrationGuideURL + ") before switching."
