package person_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	auth0ClientID     string
	auth0ClientSecret string
	auth0Audience     string
	auth0Endpoint     string
	auth0Scopes       []string
	httpClient        *http.Client
	personEndpoint    string

	auth0AccessToken string
}

func NewClient(auth0ClientID string, auth0ClientSecret string, auth0Audience string, auth0Endpoint string, auth0Scopes []string, personEndpoint string) *Client {
	c := &Client{
		auth0ClientID:     auth0ClientID,
		auth0ClientSecret: auth0ClientSecret,
		auth0Audience:     auth0Audience,
		auth0Endpoint:     auth0Endpoint,
		auth0Scopes:       auth0Scopes,
		httpClient:        &http.Client{},
		personEndpoint:    personEndpoint,
	}

	return c
}

func (client *Client) GetAccessToken(ctx context.Context) error {
	oauth2_config := clientcredentials.Config{
		ClientID:       client.auth0ClientID,
		ClientSecret:   client.auth0ClientSecret,
		EndpointParams: url.Values{"audience": {client.auth0Audience}},
		Scopes:         client.auth0Scopes,
		TokenURL:       client.auth0Endpoint,
	}

	oauth_token, err := oauth2_config.Token(ctx)
	tflog.Info(ctx, fmt.Sprintf("HTTP Request: %#v", oauth_token))

	if err == nil {
		client.auth0AccessToken = oauth_token.AccessToken
	}

	return err
}
/*
WORKAROUND: hit the dino-park whoami endpoint to get Github username from ID since cis has stale data (MZCLD-3067) 
https://github.com/mozilla-iam/dino-park-whoami/blob/377130f75a69efc52a28de8b88ca075b6dbdca9b/src/github/app.rs#L62
TODO: Fix this by doing one of:
 - Move this lookup upstream to fix cis staleness https://github.com/mozilla-iam/cis/blob/master/python-modules/cis_publisher/cis_publisher/auth0.py#L311
 - Hit Github directly instead of whoami
 - Hope https://github.com/integrations/terraform-provider-github/pull/3436 gets merged, do this lookup in TF
*/
func (client *Client) GetGithubUsernameByNodeID(ctx context.Context, githubIDV3 string) (string, error) {
	httpReq, err := http.NewRequest("GET", "https://people.mozilla.org/whoami/github/username/"+githubIDV3, nil)
	if err != nil {
		return "", err
	}

	httpResp, err := client.httpClient.Do(httpReq)
	tflog.Info(ctx, fmt.Sprintf("HTTP Request: %#v", httpReq))
	if err != nil {
		return "", err
	}

	if httpResp.StatusCode != 200 {
		return "", fmt.Errorf("people.mozilla.org responded with status code %d", httpResp.StatusCode)
	}

	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]string
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse GitHub username response: %w", err)
	}

	username, ok := result["username"]
	if !ok {
		return "", fmt.Errorf("GitHub username response missing 'username' key")
	}

	return username, nil
}

func (client *Client) GetPersonByEmail(ctx context.Context, email string) (*Person, error) {
	person := Person{}

	httpReq, err := http.NewRequest("GET", client.personEndpoint+"/v2/user/primary_email/"+email, nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Add("Authorization", "Bearer "+client.auth0AccessToken)

	httpResp, err := client.httpClient.Do(httpReq)
	tflog.Info(ctx, fmt.Sprintf("HTTP Request: %#v", httpReq))
	if err != nil {
		return nil, err
	}

	if httpResp.StatusCode >= 400 {
		return nil, fmt.Errorf("Person API responded with status code %d", httpResp.StatusCode)
	}

	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(respBody, &person)
	if err != nil {
		return nil, err
	}

	// Convert map keys into a list of strings
	keys := make([]string, 0, len(person.AccessInformation.Mozilliansorg.Values))
	for key := range person.AccessInformation.Mozilliansorg.Values {
		keys = append(keys, key)
	}
	person.AccessInformation.Mozilliansorg.List = keys

	return &person, nil
}
