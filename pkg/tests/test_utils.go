package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/glide/sdk-go/pkg/types"
	"github.com/joho/godotenv"
)

func getEnvOrDefault(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

func SetupTestEnvironment(t *testing.T) types.GlideSdkSettings {
    // Try to load .env file from multiple possible locations
    envLocations := []string{
        ".env",
        "../.env",
        "../../.env",
    }

    var loadedEnv bool
    for _, loc := range envLocations {
        err := godotenv.Load(loc)
        if err == nil {
            t.Logf("Loaded .env file from: %s", loc)
            loadedEnv = true
            break
        }
    }

    if !loadedEnv {
        t.Logf("Warning: Failed to load .env file from any location. Using environment variables directly.")
    }

    // Create settings from environment variables
    return types.GlideSdkSettings{
        ClientID:     os.Getenv("GLIDE_CLIENT_ID"),
        ClientSecret: os.Getenv("GLIDE_CLIENT_SECRET"),
        RedirectURI:  os.Getenv("GLIDE_REDIRECT_URI"),
        Internal: types.InternalSettings{
            AuthBaseURL: getEnvOrDefault("GLIDE_AUTH_BASE_URL", "https://oidc.gateway-x.io"),
            APIBaseURL:  getEnvOrDefault("GLIDE_API_BASE_URL", "https://api.gateway-x.io"),
        },
    }
}

type HttpResponse struct {
	Headers http.Header
	Data  string
	Query url.Values
}

func MakeRawHttpRequestFollowRedirectChain(urlStr string) (*HttpResponse, error) {
	// Create a cookie jar
	fmt.Println("urlStr "+urlStr)
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}
	// Create a client with the cookie jar
	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 10 redirects
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			fmt.Println("Redirecting to: ", req.URL)
			if strings.HasPrefix(req.URL.String(), "https://playground.glideapi.com/magical-auth/verify") {
				fmt.Println("stopping redirect")
				return http.ErrUseLastResponse
			} else if strings.HasPrefix(req.URL.String(), "https://dev.gateway-x.io/dev-redirector/callback?code") {
				fmt.Println("stopping redirect")
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	// Parse the initial URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}
	// Make the request
	resp, err := client.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Determine the data to return
	var data string
	if token := resp.Header.Get("Token"); token != "" {
		data = token
	} else {
		data = string(body)
	}

	// Create and return the HttpResponse
	return &HttpResponse{
		Data:  data,
		Query: parsedURL.Query(),
		Headers: resp.Header,
	}, nil
}

// func MakeRawHttpRequestFollowRedirectChain(urlStr string) (HttpResponse, error) {
// 	client := &http.Client{
// 		CheckRedirect: func(req *http.Request, via []*http.Request) error {
// 			return http.ErrUseLastResponse
// 		},
// 	}
// 	for {
// 		fmt.Println("urlStr "+urlStr)
// 		resp, err := client.Get(urlStr)
// 		fmt.Println("resp")
// 		fmt.Println(resp)
// 		if err != nil {
// 			return HttpResponse{}, err
// 		}
// 		defer resp.Body.Close()
//
// 		if resp.StatusCode == http.StatusOK {
// 			body, err := ioutil.ReadAll(resp.Body)
// 			if err != nil {
// 				return HttpResponse{}, err
// 			}
//
// 			data := string(body)
// 			if token := resp.Header.Get("Token"); token != "" {
// 				data = token
// 			}
//
// 			parsedURL, err := url.Parse(urlStr)
// 			if err != nil {
// 				return HttpResponse{}, err
// 			}
//
// 			return HttpResponse{
// 				Data:  data,
// 				Query: parsedURL.Query(),
// 			}, nil
// 		}
//
// 		if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther || resp.StatusCode == http.StatusTemporaryRedirect {
// 			location := resp.Header.Get("Location")
// 			if location == "" {
// 				return HttpResponse{}, fmt.Errorf("redirect with no location")
// 			}
// 			urlStr = location
// 			continue
// 		}
// 		return HttpResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
// 	}
// }
