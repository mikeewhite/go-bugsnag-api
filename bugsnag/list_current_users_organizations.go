package bugsnag

import (
	"context"
	"time"
)

type Organization struct {
	AutoUpgrade      bool      `json:"auto_upgrade"`
	BillingEmails    []string  `json:"billing_emails"`
	CollaboratorsURL string    `json:"collaborators_url"`
	CreatedAt        time.Time `json:"created_at"`
	Creator          Creator   `json:"creator"`
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	ProjectsURL      string    `json:"projects_url"`
	Slug             string    `json:"slug"`
	UpdatedAt        time.Time `json:"updated_at"`
	UpgradeURL       string    `json:"upgrade_url"`
}

type Creator struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

type ListCurrentUsersOrganizationsOptions struct {
	// Whether all organizations should be returned (false) or only
	// organizations that the current user is an admin of (true)
	Admin *bool `url:"admin,omitempty"`

	// The maximum number of results to return in each page of results
	PerPage *int `url:"per_page,omitempty"`
}

// ListCurrentUsersOrganizations lists the organizations that the 'current user' is a member of
//
// Bugsnag API docs: https://bugsnagapiv2.docs.apiary.io/#reference/current-user/organizations/list-the-current-user's-organizations
func (c *Client) ListCurrentUsersOrganizations(ctx context.Context, opts *ListCurrentUsersOrganizationsOptions) ([]*Organization, *Response, error) {
	endpoint, err := addOptions("user/organizations", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, nil, err
	}

	orgs := new([]*Organization)
	resp, err := c.Do(ctx, req, orgs)

	return *orgs, resp, err
}
