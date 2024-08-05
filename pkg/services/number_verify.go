package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
	"github.com/google/uuid"
	"github.com/glide/sdk-go/pkg/types"
	"github.com/glide/sdk-go/pkg/utils"
)

type NumberVerifyAuthUrlInput struct {
    State *string
}


type NumberVerifyFuncResponse struct {
	DevicePhoneNumberVerified bool
}

type NumberVerifyUserClient struct {
    settings types.GlideSdkSettings
	session  *types.Session
	code        string
	phoneNumber *string
}

func NewNumberVerifyUserClient(settings types.GlideSdkSettings, params types.NumberVerifyClientForParams) *NumberVerifyUserClient {
	return &NumberVerifyUserClient{
		settings:    settings,
		code:        params.Code,
		phoneNumber: params.PhoneNumber,
	}
}

func (c *NumberVerifyUserClient) StartSession() error {
	if c.settings.Internal.AuthBaseURL == "" {
		return errors.New("[GlideClient] internal.authBaseUrl is unset")
	}
	if c.settings.ClientID == "" || c.settings.ClientSecret == "" {
		return errors.New("[GlideClient] Client credentials are required to generate a new session")
	}
	if c.code == "" {
		return errors.New("[GlideClient] Code is required to start a session")
	}
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", c.code)
	// TODO: Set the redirect_uri to the actual redirect uri
	data.Set("redirect_uri", "https://dev.gateway-x.io/dev-redirector/callback")
    resp, err := utils.FetchX(c.settings.Internal.AuthBaseURL+"/oauth2/token", utils.FetchXInput{
        Method: "POST",
        Headers: map[string]string{
            "Content-Type":  "application/x-www-form-urlencoded",
            "Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(c.settings.ClientID+":"+c.settings.ClientSecret)),
        },
        Body: data.Encode(),
    })
	if err != nil {
		return fmt.Errorf("failed to generate new session: %w", err)
	}
	var body struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64    `json:"expires_in"`
		Scope       string `json:"scope"`
	}

	if err := resp.JSON(&body); err != nil {
	        fmt.Errorf("[GlideClient] Failed to parse response: %w", err)
    		return nil
    }

	c.session = &types.Session{
		AccessToken: body.AccessToken,
		ExpiresAt:   time.Now().Unix() + body.ExpiresIn,
		Scopes:      strings.Split(body.Scope, " "),
	}
	return nil
}

func (c *NumberVerifyUserClient) VerifyNumber(number ...string) (*NumberVerifyFuncResponse, error) {
	if c.session == nil {
		return nil, errors.New("[GlideClient] Session is required to verify a number")
	}

	if c.settings.Internal.APIBaseURL == "" {
		return nil, errors.New("[GlideClient] internal.apiBaseUrl is unset")
	}

	var phoneNumber string
	if len(number) > 0 && number[0] != "" {
		phoneNumber = number[0]
	} else if c.phoneNumber != nil {
		phoneNumber = *c.phoneNumber
	}

	// phoneNumber := number
	// if phoneNumber == "" {
	// 	phoneNumber = c.phoneNumber
	// }
	if phoneNumber == "" {
		return nil, errors.New("[GlideClient] Phone number is required to verify a number")
	}

	body, err := json.Marshal(map[string]string{"phoneNumber": utils.FormatPhoneNumber(phoneNumber)})
	if err != nil {
		return nil, fmt.Errorf("[GlideClient] failed to marshal payload in number verify: %w", err)
	}

	resp, err := utils.FetchX(c.settings.Internal.APIBaseURL+"/number-verification/verify", utils.FetchXInput{
    		Method: "POST",
    		Headers: map[string]string{
    			"Content-Type":  "application/json",
    			"Authorization": "Bearer " + c.session.AccessToken,
    		},
    		Body: string(body),
    })

	if err != nil {
		return nil, fmt.Errorf("failed to verify number: %w", err)
	}

	// You might want to process the response body here
    fmt.Println(resp)
	var result NumberVerifyFuncResponse
	if err := resp.JSON(&result); err != nil {
		return nil, fmt.Errorf("[GlideClient] Failed to parse response: %w", err)
	}
	return &result, nil
}

type NumberVerifyClient struct {
	settings types.GlideSdkSettings
}

func NewNumberVerifyClient(settings types.GlideSdkSettings) *NumberVerifyClient {
	return &NumberVerifyClient{settings: settings}
}

func (c *NumberVerifyClient) GetAuthURL(opts ...NumberVerifyAuthUrlInput) (string, error) {
	if c.settings.Internal.AuthBaseURL == "" {
		return "", errors.New("[GlideClient] internal.authBaseUrl is unset")
	}
	if c.settings.ClientID == "" {
		return "", errors.New("[GlideClient] Client id is required to generate an auth url")
	}
	var state string
    if len(opts) > 0 && opts[0].State != nil {
        state = *opts[0].State
    } else {
        state = uuid.New().String()
    }
	nonce := uuid.New().String()
	params := url.Values{}
	params.Set("client_id", c.settings.ClientID)
	params.Set("response_type", "code")
	if c.settings.RedirectURI != "" {
		params.Set("redirect_uri", c.settings.RedirectURI)
	}
	params.Set("scope", "openid")
	params.Set("purpose", "dpv:FraudPreventionAndDetection:number-verification")
	params.Set("state", state)
	params.Set("nonce", nonce)
	params.Set("max_age", "0")

	return c.settings.Internal.AuthBaseURL + "/oauth2/auth?" + params.Encode(), nil
}

func (c *NumberVerifyClient) For(params types.NumberVerifyClientForParams) (*NumberVerifyUserClient, error) {
	client := NewNumberVerifyUserClient(c.settings, params)
	err := client.StartSession()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *NumberVerifyClient) GetHello() (string) {
	return "Hello"
}

