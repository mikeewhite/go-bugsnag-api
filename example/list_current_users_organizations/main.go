package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mikeewhite/go-bugsnag-api/bugsnag"
)

func main() {
	token := os.Getenv("BUGSNAG_AUTH_TOKEN")
	ctx := context.Background()
	client, err := bugsnag.NewClient(bugsnag.WithAuthenticationToken(token))
	if err != nil {
		log.Fatalf("Error on initializing client: %s", err)
	}

	opts := &bugsnag.ListCurrentUsersOrganizationsOptions{
		Admin: bugsnag.Bool(true),
	}
	orgs, _, err := client.ListCurrentUsersOrganizations(ctx, opts)
	if err != nil {
		log.Fatalf("Error on calling ListCurrentUsersOrganizations endpoint: %s", err)
	}
	fmt.Printf("Returned %d organization(s):\n", len(orgs))
	for _, org := range orgs {
		fmt.Printf("%s\n", prettyPrint(org))
	}
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
