package tests

import (
	"net/url"
	"testing"
	"github.com/glide/sdk-go/pkg/glide"
	"github.com/glide/sdk-go/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestMagicAuth(t *testing.T) {
	settings := SetupTestEnvironment(t)
	glideClient, err := glide.NewGlideClient(settings)
	assert.NoError(t, err)
	t.Run("should start magic auth", func(t *testing.T) {
		magicRes, err := glideClient.MagicAuth.StartAuth(types.MagicAuthStartProps{
			PhoneNumber: "+555123456789",
		}, types.ApiConfig{})
		assert.NoError(t, err)
		assert.NotNil(t, magicRes)
		assert.Equal(t, "MAGIC", magicRes.Type)
		assert.NotEmpty(t, magicRes.AuthURL)
		t.Logf("Magic auth StartAuth response: %+v", magicRes)
		res, _ := MakeRawHttpRequestFollowRedirectChain(magicRes.AuthURL)
		location := res.Headers.Get("Location")
		parsedLocation, err := url.Parse(location)
		assert.NoError(t, err)
        token := parsedLocation.Query().Get("token")
		t.Logf("Magic auth response: %+v", res)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		verify, err := glideClient.MagicAuth.VerifyAuth(types.MagicAuthVerifyProps{
			PhoneNumber: "+555123456789",
			Token:       token,
		}, types.ApiConfig{})
		assert.NoError(t, err)
		t.Logf("Check verify: %+v", verify)
		assert.True(t, verify)
	})
}
