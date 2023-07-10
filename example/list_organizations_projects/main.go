package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mikeewhite/go-bugsnag-api/bugsnag"
)

var (
	orgID = flag.String("orgID", "", "ID of the organization to get projects for")
)

func main() {
	token := os.Getenv("BUGSNAG_AUTH_TOKEN")
	if token == "" {
		log.Fatalf("Authorization token should be provided as BUGSNAG_AUTH_TOKEN env var")
	}

	flag.Parse()
	if *orgID == "" {
		log.Fatalf("Organization ID should be provided via the -orgID flag")
	}

	ctx := context.Background()
	client, err := bugsnag.NewClient(bugsnag.WithAuthenticationToken(token))
	if err != nil {
		log.Fatalf("Error on initializing client: %s", err)
	}

	opts := &bugsnag.ListOrganizationsProjectsOptions{
		PerPage: bugsnag.Int(50),
	}

	pageCount := 1
	for {
		projects, resp, err := client.ListOrganizationsProjects(ctx, *orgID, opts)
		if err != nil {
			log.Fatalf("Error on calling ListOrganizationsProjects endpoint: %s", err)
		}
		fmt.Printf("Returned %d project(s) [page=%d, total=%d]:\n", len(projects), pageCount, *resp.TotalCount)
		for _, p := range projects {
			fmt.Printf("%s\n", prettyPrint(p))
		}
		if resp.NextPageURL == nil {
			break
		}
		opts.NextPageURL = resp.NextPageURL
		pageCount++
	}
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
