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

func TestListCurrentUsersOrganizations(t *testing.T) {
	tv := setup(t)

	tv.mux.HandleFunc("/user/organizations", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		json, err := os.ReadFile("testdata/list_current_users_organizations_response.json")
		require.NoError(t, err)
		fmt.Fprint(w, string(json))
		w.WriteHeader(http.StatusOK)
	})

	opts := &ListCurrentUsersOrganizationsOptions{}
	ctx := context.Background()
	orgs, _, err := tv.client.ListCurrentUsersOrganizations(ctx, opts)

	require.NoError(t, err)
	require.NotNil(t, orgs)
	require.NotEmpty(t, orgs)
	assert.Len(t, orgs, 1)

	org := orgs[0]
	assert.Equal(t, "515fb9337c1074f6fd000007", org.ID)
	assert.Equal(t, "Acme Co.", org.Name)
	assert.Equal(t, "acme-co", org.Slug)
	require.NotNil(t, org.Creator)
	assert.Equal(t, "user@example.com", org.Creator.Email)
	assert.Equal(t, "58c9b9b09ef17217f1fb8b30", org.Creator.ID)
	assert.Equal(t, "Joe Bloggs", org.Creator.Name)
	assert.Equal(t, "https://api.bugsnag.com/organizations/515fb9337c1074f6fd000007/collaborators", org.CollaboratorsURL)
	assert.Equal(t, "https://api.bugsnag.com/organizations/515fb9337c1074f6fd000007/projects", org.ProjectsURL)
	assert.Equal(t, "2017-04-24T22:17:13Z", org.CreatedAt.Format(time.RFC3339))
	assert.Equal(t, "2017-04-24T22:17:13Z", org.UpdatedAt.Format(time.RFC3339))
	assert.True(t, org.AutoUpgrade)
	assert.Equal(t, "https://api.bugsnag.com/settings/bugsnag/plans-billing?plansBilling%5Bstep%5D=collaborators-and-events", org.UpgradeURL)
	require.Len(t, org.BillingEmails, 1)
	assert.Equal(t, "user@example.com", org.BillingEmails[0])
}
