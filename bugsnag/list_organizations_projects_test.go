package bugsnag

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListOrganizationsProjects(t *testing.T) {
	tv := setup(t)
	orgID := "547c4b0b69196200109ead5c"
	url := fmt.Sprintf("/organizations/%s/projects", orgID)

	tv.mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		json, err := os.ReadFile("testdata/list_organizations_projects_response.json")
		require.NoError(t, err)
		fmt.Fprint(w, string(json))
		w.WriteHeader(http.StatusOK)
	})

	opts := &ListOrganizationsProjectsOptions{}
	ctx := context.Background()
	projs, _, err := tv.client.ListOrganizationsProjects(ctx, orgID, opts)

	require.NoError(t, err)
	require.NotNil(t, projs)
	require.NotEmpty(t, projs)
	assert.Len(t, projs, 1)

	p := projs[0]
	assert.Equal(t, "537c4b0b69196200109eac5c", p.ID)
	assert.Equal(t, orgID, p.OrganizationID)
	assert.Equal(t, "example-project", p.Slug)
	assert.Equal(t, "Example Project", p.Name)
	assert.Equal(t, "2e3b1d5af480d995d80d1536442117d5", p.APIKey)
	assert.Equal(t, "react", p.Type)
	assert.True(t, p.IsFullView)
	assert.Equal(t, []string{"staging", "production"}, p.ReleaseStages)
	assert.Equal(t, "javascript", p.Language)
	assert.Equal(t, "2017-04-24T22:17:13Z", p.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, "2017-04-24T22:17:13Z", p.UpdatedAt.Format(time.RFC3339))
	assert.Equal(t, "https://api.bugsnag.com/projects/537c4b0b69196200109eac5c/errors", p.ErrorsURL)
	assert.Equal(t, "https://api.bugsnag.com/projects/537c4b0b69196200109eac5c/events", p.EventsURL)
	assert.Equal(t, "https://api.bugsnag.com/projects/537c4b0b69196200109eac5c", p.URL)
	assert.Equal(t, "https://app.bugsnag.com/example-account/example-project", p.HTMLURL)
	assert.Equal(t, 1, p.OpenErrorCount)
	assert.Equal(t, 2, p.ForReviewErrorCount)
	assert.Equal(t, 49, p.CollaboratorsCount)
	assert.Equal(t, 1, p.TeamsCount)
	assert.Empty(t, p.GlobalGrouping)
	assert.Empty(t, p.LocationGrouping)
	assert.Empty(t, p.DiscardedAppVersions)
	assert.Empty(t, p.DiscardedErrors)
	assert.Equal(t, 5, p.CustomEventFieldsUsed)
	assert.False(t, p.ResolveOnDeploy)
	assert.Equal(t, []string{"example.com"}, p.URLWhitelist)
	assert.True(t, p.IgnoreOldBrowsers)
	assert.Empty(t, p.IgnoredBrowserVersions)
}
