package glide

import (
	"errors"
	"os"

	"github.com/glide/sdk-go/pkg/services"
	"github.com/glide/sdk-go/pkg/types"
)

// GlideClient is the main client for the SDK
type GlideClient struct {
	settings    types.GlideSdkSettings
	TelcoFinder *services.TelcoFinderClient
	MagicAuth   *services.MagicAuthClient
	SimSwap     *services.SimSwapClient
	NumberVerify *services.NumberVerifyClient
}

func NewGlideClient(settings types.GlideSdkSettings) (*GlideClient, error) {
	defaults := types.GlideSdkSettings{
		ClientID:     os.Getenv("GLIDE_CLIENT_ID"),
		ClientSecret: os.Getenv("GLIDE_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("GLIDE_REDIRECT_URI"),
		Internal: types.InternalSettings{
			AuthBaseURL: getEnvOrDefault("GLIDE_AUTH_BASE_URL", "https://oidc.gateway-x.io"),
			APIBaseURL:  getEnvOrDefault("GLIDE_API_BASE_URL", "https://api.gateway-x.io"),
		},
	}

	// Merge defaults with provided settings
	mergedSettings := mergeSettings(defaults, settings)

	if mergedSettings.ClientID == "" {
		return nil, errors.New("clientId is required")
	}

	if mergedSettings.Internal.AuthBaseURL == "" {
		return nil, errors.New("internal.authBaseUrl is unset")
	}

	client := &GlideClient{
		settings:    mergedSettings,
		TelcoFinder: services.NewTelcoFinderClient(mergedSettings),
		MagicAuth:   services.NewMagicAuthClient(mergedSettings),
		SimSwap:     services.NewSimSwapClient(mergedSettings),
		NumberVerify: services.NewNumberVerifyClient(mergedSettings),
	}

	return client, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func mergeSettings(defaults, override types.GlideSdkSettings) types.GlideSdkSettings {
	result := defaults

	if override.ClientID != "" {
		result.ClientID = override.ClientID
	}
	if override.ClientSecret != "" {
		result.ClientSecret = override.ClientSecret
	}
	if override.RedirectURI != "" {
		result.RedirectURI = override.RedirectURI
	}
	if override.Internal.AuthBaseURL != "" {
		result.Internal.AuthBaseURL = override.Internal.AuthBaseURL
	}
	if override.Internal.APIBaseURL != "" {
		result.Internal.APIBaseURL = override.Internal.APIBaseURL
	}

	return result
}