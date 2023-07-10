package bugsnag

import (
	"context"
	"fmt"
	"net/url"
	"time"
)

type Project struct {
	Name                   string                 `json:"name"`
	GlobalGrouping         []string               `json:"global_grouping"`
	LocationGrouping       []string               `json:"location_grouping"`
	DiscardedAppVersions   []string               `json:"discarded_app_versions"`
	DiscardedErrors        []string               `json:"discarded_errors"`
	URLWhitelist           []string               `json:"url_whitelist"`
	IgnoreOldBrowsers      bool                   `json:"ignore_old_browsers"`
	IgnoredBrowserVersions map[string]interface{} `json:"ignored_browser_versions"`
	ResolveOnDeploy        bool                   `json:"resolve_on_deploy"`
	ID                     string                 `json:"id"`
	OrganizationID         string                 `json:"organization_id"`
	Type                   string                 `json:"type"`
	Slug                   string                 `json:"slug"`
	APIKey                 string                 `json:"api_key"`
	IsFullView             bool                   `json:"is_full_view"`
	ReleaseStages          []string               `json:"release_stages"`
	Language               string                 `json:"language"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`
	URL                    string                 `json:"url"`
	HTMLURL                string                 `json:"html_url"`
	ErrorsURL              string                 `json:"errors_url"`
	EventsURL              string                 `json:"events_url"`
	OpenErrorCount         int                    `json:"open_error_count"`
	ForReviewErrorCount    int                    `json:"for_review_error_count"`
	CollaboratorsCount     int                    `json:"collaborators_count"`
	TeamsCount             int                    `json:"teams_count"`
	CustomEventFieldsUsed  int                    `json:"custom_event_fields_used"`
}

type ListOrganizationsProjectsOptions struct {
	// The maximum number of results to return in each page of results
	PerPage *int `url:"per_page,omitempty"`

	NextPageURL *url.URL `url:"-"`
}

// ListOrganizationsProjects lists the projects associated with a given organization
//
// See https://bugsnagapiv2.docs.apiary.io/#reference/current-user/organizations/list-an-organization's-projects
func (c *Client) ListOrganizationsProjects(ctx context.Context, orgID string, opts *ListOrganizationsProjectsOptions) ([]*Project, *Response, error) {
	var endpoint string
	if opts.NextPageURL != nil {
		endpoint = fmt.Sprintf("%s?%s", opts.NextPageURL.Path[1:], opts.NextPageURL.RawQuery)
	} else {
		var err error
		endpoint, err = addOptions(fmt.Sprintf("organizations/%s/projects", orgID), opts)
		if err != nil {
			return nil, nil, err
		}
	}

	req, err := c.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	projs := new([]*Project)
	resp, err := c.Do(ctx, req, projs)

	return *projs, resp, err
}
