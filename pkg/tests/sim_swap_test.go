package tests

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/glide/sdk-go/pkg/types"
    "github.com/glide/sdk-go/pkg/glide"
)

func TestSimSwapClient(t *testing.T) {
	settings := SetupTestEnvironment(t)
	client, err := glide.NewGlideClient(settings)
    assert.NoError(t, err, "NewGlideClient should not return an error")
    t.Run("For", func(t *testing.T) {
        userClient, err := client.SimSwap.For(types.PhoneIdentifier{PhoneNumber: "+555123456789"})
        assert.NoError(t, err, "For should not return an error")
        assert.NotNil(t, userClient, "UserClient should not be nil")
    })

    t.Run("Check", func(t *testing.T) {
        userClient, _ := client.SimSwap.For(types.PhoneIdentifier{PhoneNumber: "+555123456789"})
        response, err := userClient.Check(types.SimSwapCheckParams{PhoneNumber: "+555123456789"}, types.ApiConfig{})
        assert.NoError(t, err, "Check should not return an error")
        assert.NotNil(t, response, "Response should not be nil")
        assert.Equal(t, true, response.Swapped, "Response should have swapped=true")
        t.Logf("Check response: %+v", response)
    })

    t.Run("RetrieveDate", func(t *testing.T) {
        userClient, _ := client.SimSwap.For(types.PhoneIdentifier{PhoneNumber: "+555123456789"})
        response, err := userClient.RetrieveDate(types.SimSwapRetrieveDateParams{PhoneNumber: "+555123456789"}, types.ApiConfig{})
        assert.NoError(t, err, "RetrieveDate should not return an error")
        assert.NotNil(t, response, "Response should not be nil")
        t.Logf("RetrieveDate response: %+v", response)
        // Add more specific assertions based on the expected response
    })

    // t.Run("GetConsentURL", func(t *testing.T) {
    //     userClient, _ := client.SimSwap.For(types.PhoneIdentifier{PhoneNumber: "+555123456789"})
    //     consentURL := userClient.GetConsentURL()
    //     assert.NotEmpty(t, consentURL, "ConsentURL should not be empty")
    // })

    t.Run("PollAndWaitForSession", func(t *testing.T) {
        userClient, _ := client.SimSwap.For(types.PhoneIdentifier{PhoneNumber: "+555123456789"})
        err := userClient.PollAndWaitForSession()
        assert.NoError(t, err, "PollAndWaitForSession should not return an error")
    })
}
