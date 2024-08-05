package tests

import (
	"net/url"
	"testing"
	"github.com/glide/sdk-go/pkg/glide"
	"github.com/glide/sdk-go/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestNumberVerify(t *testing.T) {
	settings := SetupTestEnvironment(t)
	glideClient, err := glide.NewGlideClient(settings)
	assert.NoError(t, err)
	t.Run("should work", func(t *testing.T) {
	    authUrl, err := glideClient.NumberVerify.GetAuthURL();
		assert.NoError(t, err)
		assert.NotNil(t, authUrl)
		assert.NotEmpty(t, authUrl)
		t.Logf("authUrl response: %+v", authUrl)
		baseURL, err := url.Parse(authUrl)
		if err != nil {
			t.Fatalf("Failed to parse authUrl: %v", err)
		}
		query := baseURL.Query()
		query.Set("login_hint", "tel:+555123456789")
		baseURL.RawQuery = query.Encode()
		res, _ := MakeRawHttpRequestFollowRedirectChain(baseURL.String())
		location := res.Headers.Get("Location")
		parsedLocation, err := url.Parse(location)
		assert.NoError(t, err)
        code := parsedLocation.Query().Get("code")
		t.Logf("Code: %s", code)
		assert.NoError(t, err)
		assert.NotEmpty(t, code)
		phoneNumber := "+555123456789"
		client, err := glideClient.NumberVerify.For(types.NumberVerifyClientForParams{PhoneNumber: &phoneNumber, Code: code})
		assert.NoError(t, err)
		verify, err := client.VerifyNumber()
		assert.NoError(t, err)
		t.Logf("Check verify: %+v", verify)
		assert.Equal(t, true, verify.DevicePhoneNumberVerified, "Response should have devicePhoneNumberVerified=true")
	})
}
