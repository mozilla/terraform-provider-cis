package provider

import (
	"runtime"
	"strings"
	"testing"
)

func TestBuildUserAgent(t *testing.T) {
	platform := runtime.GOOS + "/" + runtime.GOARCH
	goVersion := runtime.Version()

	cases := []struct {
		name             string
		providerVersion  string
		terraformVersion string
		want             string
	}{
		{
			name:             "release build",
			providerVersion:  "1.2.3",
			terraformVersion: "1.9.0",
			want:             "terraform-provider-cis/1.2.3 (release) Terraform/1.9.0 Go/" + goVersion + " (" + platform + ")",
		},
		{
			name:             "dev build",
			providerVersion:  "dev",
			terraformVersion: "1.9.0",
			want:             "terraform-provider-cis/dev (dev) Terraform/1.9.0 Go/" + goVersion + " (" + platform + ")",
		},
		{
			name:             "test build",
			providerVersion:  "test",
			terraformVersion: "1.6.0",
			want:             "terraform-provider-cis/test (dev) Terraform/1.6.0 Go/" + goVersion + " (" + platform + ")",
		},
		{
			name:             "empty versions",
			providerVersion:  "",
			terraformVersion: "",
			want:             "terraform-provider-cis/unknown (dev) Terraform/unknown Go/" + goVersion + " (" + platform + ")",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildUserAgent(tc.providerVersion, tc.terraformVersion)
			if got != tc.want {
				t.Errorf("buildUserAgent(%q, %q) = %q, want %q", tc.providerVersion, tc.terraformVersion, got, tc.want)
			}
			if !strings.HasPrefix(got, "terraform-provider-cis/") {
				t.Errorf("user agent %q missing product prefix", got)
			}
		})
	}
}
