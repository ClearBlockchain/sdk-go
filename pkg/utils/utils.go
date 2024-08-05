package utils

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"
    "strings"
)

// HTTPResponseError represents an HTTP error response
type HTTPResponseError struct {
    Response *http.Response
}

func (e *HTTPResponseError) Error() string {
    return fmt.Sprintf("HTTP Error Response: %d %s", e.Response.StatusCode, e.Response.Status)
}

// InsufficientSessionError represents an error due to insufficient session
type InsufficientSessionError struct {
    Have    *int
    Need    *int
    Message string
}

func (e *InsufficientSessionError) Error() string {
    if e.Message != "" {
        return e.Message
    }
    return "Session is required for this request"
}

// FormatPhoneNumber formats a phone number string
func FormatPhoneNumber(phoneNumber string) string {
    re := regexp.MustCompile("[^0-9]")
    return "+" + re.ReplaceAllString(phoneNumber, "")
}

// FetchError represents an error during fetch operation
type FetchError struct {
    Response *http.Response
    Data     string
}

func (e *FetchError) Error() string {
    return fmt.Sprintf("Fetch Error: %d %s", e.Response.StatusCode, e.Response.Status)
}

// FetchXInput represents input for FetchX function
type FetchXInput struct {
    Method  string
    Headers map[string]string
    Body    string
}

// FetchXResponse represents the response from FetchX function
type FetchXResponse struct {
    Data []byte
    Response *http.Response
}

func (r *FetchXResponse) JSON(v interface{}) error {
    return json.Unmarshal(r.Data, v)
}

func (r *FetchXResponse) Text() string {
    return string(r.Data)
}

func (r *FetchXResponse) OK() bool {
    return r.Response.StatusCode < 400
}

// FetchX performs an HTTP request
func FetchX(url string, input FetchXInput) (*FetchXResponse, error) {
    client := &http.Client{}

    req, err := http.NewRequest(input.Method, url, strings.NewReader(input.Body))
    if err != nil {
        return nil, err
    }

    for k, v := range input.Headers {
        req.Header.Set(k, v)
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode >= 400 {
        return nil, &FetchError{Response: resp, Data: string(data)}
    }

    return &FetchXResponse{Data: data, Response: resp}, nil
}
