package bugsnag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	defaultAPIVersion = "2"
	defaultBaseURL    = "https://api.bugsnag.com/"

	headerAPIVersion = "X-Version"
)

// Option is a functional option for configuring the API client
type Option func(*Client) error

// Client for interacting with the Bugsnag data access API
// See https://bugsnagapiv2.docs.apiary.io/
type Client struct {
	// The access token for the data access API.
	// See https://bugsnagapiv2.docs.apiary.io/#introduction/authentication/personal-auth-tokens-(recommended)
	authenticationToken string

	// The base URL for API requests. Defaults to the public Bugsnag API, but can be
	// overridden for use with on-premise installations.
	baseURL *url.URL

	httpClient *http.Client
}

// Response is a Bugsnag API response. This wraps the standard http.Response returned from
// Bugsnag and provides pagination controls
type Response struct {
	*http.Response

	// Most requests that return multiple items will also include the total count of all results
	TotalCount *int64

	NextPageURL *url.URL
}

// NewClient creates a new client for interacting with the Bugsnag data access API
func NewClient(opts ...Option) (*Client, error) {
	baseURL, _ := url.Parse(defaultBaseURL)
	client := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}

	// apply any options
	for _, option := range opts {
		err := option(client)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

// WithAuthenticationToken sets the provided authentication token
// See https://bugsnagapiv2.docs.apiary.io/#introduction/authentication/personal-auth-tokens-(recommended)
func WithAuthenticationToken(token string) Option {
	return func(c *Client) error {
		c.authenticationToken = token
		return nil
	}
}

// WithBaseURL allows the base URL of the API client to be overridden.
// The provided URL must end with a trailing slash
func WithBaseURL(baseURLStr string) Option {
	return func(c *Client) error {
		baseURL, err := url.Parse(baseURLStr)
		if err != nil {
			return err
		}
		if !strings.HasSuffix(c.baseURL.Path, "/") {
			return fmt.Errorf("baseURL must have a trailing slash")
		}
		c.baseURL = baseURL

		return nil
	}
}

// NewRequest creates an API request
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set(headerAPIVersion, defaultAPIVersion)
	if c.authenticationToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", c.authenticationToken))
	}

	return req, nil
}

// Do performs a HTTP request
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// if the context was cancelled then return the error from the context
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return nil, err
		}
	}
	defer resp.Body.Close()

	response, err := parseResponse(resp)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Unexpected response code: %d", response.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(v)
	if err == io.EOF {
		err = nil // ignore EOF errors caused by empty response body
	}
	if err != nil {
		return nil, err
	}

	return response, err
}

func parseResponse(r *http.Response) (*Response, error) {
	response := &Response{Response: r}

	// parse the 'total result count' header if present
	// See https://bugsnagapiv2.docs.apiary.io/#introduction/pagination/total-count-of-results
	totalResults := r.Header.Get("X-Total-Count")
	if totalResults != "" {
		count, err := strconv.Atoi(totalResults)
		if err != nil {
			return nil, err
		}
		response.TotalCount = Int64(int64(count))
	}

	// parse the link header (used for pagination)
	// See https://bugsnagapiv2.docs.apiary.io/#introduction/pagination
	if links, ok := r.Header["Link"]; ok && len(links) > 0 {
		// if there were multiple 'link' headers then they will be concatenated with comma separators
		for _, link := range strings.Split(links[0], ",") {
			// split the link header into its constituent parts where the href and the rel should be separated by a semi-colon, e.g.
			// link: <https://api.bugsnag.com/organizations/3abaed0d9bf39c1431000001/projects?direction=desc&offset%5Bnull_sort_field%5D=false&offset%5Bsort_field_offset%5D=601ac082ec80d80015bc0a85&per_page=30&sort=created_at>; rel="next"
			segments := strings.Split(strings.TrimSpace(link), ";")

			// link must contain at least a href and rel
			if len(segments) < 2 {
				continue
			}

			// ensure href is properly formatted
			href := segments[0]
			if !strings.HasPrefix(href, "<") || !strings.HasSuffix(href, ">") {
				continue
			}

			// if this link represents the next page link then store it against the response
			if strings.TrimSpace(segments[1]) == `rel="next"` {
				// parse the URL first to ensure it is well formed
				url, err := url.Parse(href[1 : len(segments[0])-1])
				if err != nil {
					continue
				}
				response.NextPageURL = url
			}

		}
	}
	return response, nil
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// Bool is a helper function that returns a pointer to a bool value
func Bool(v bool) *bool { return &v }

// Int is a helper function that returns a pointer to an int value
func Int(v int) *int { return &v }

// Int64 is a helper function that returns a pointer to an int64 value
func Int64(v int64) *int64 { return &v }

// String is a helper function that returns a pointer to a string value
func String(v string) *string { return &v }
