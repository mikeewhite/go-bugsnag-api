# 🐛 go-bugsnag-api

go-bugsnag-api is a draft implementation of a Go client library for accessing the [Bugsnag data access API](https://bugsnagapiv2.docs.apiary.io/#) inspired by [go-github](https://github.com/google/go-github).

## Installation

```bash
go get github.com/mikeewhite/go-bugsnag-api
```

## Usage

```go
import "github.com/mikeewhite/go-bugsnag-api/bugsnag"

client := bugsnag.NewClient(nil)

// List all organizations for the 'current user'
opts := &bugsnag.ListCurrentUsersOrganizationsOptions{
     Admin: bugsnag.Bool(true),
}
orgs, _, err := client.ListCurrentUsersOrganizations(ctx, opts)
```

## TODOs
 - [ ] Implement [Errors API](https://bugsnagapiv2.docs.apiary.io/#reference/errors/list-an-organization's-projects)
 - [ ] Implement [Integrations API](https://bugsnagapiv2.docs.apiary.io/#reference/integrations/list-an-organization's-projects)
 - [ ] Implement [Organizations API](https://bugsnagapiv2.docs.apiary.io/#reference/organizations/list-an-organization's-projects)
 - [ ] Implment [Projects API](https://bugsnagapiv2.docs.apiary.io/#reference/projects/list-an-organization's-projects)
