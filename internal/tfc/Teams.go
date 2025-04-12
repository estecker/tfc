package tfc

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"log"
)

func TeamsCmd(org string) {
	fmt.Println("Using organization:", org)
	listTeams(org)
}

func listTeams(org string) {
	client, err := Client()
	ctx := context.Background()
	if err != nil {
		log.Fatalf("Issues creating a client\n")
	}
	tl, err := client.Teams.List(ctx,
		org,
		&tfe.TeamListOptions{ListOptions: tfe.ListOptions{},
			Include: []tfe.TeamIncludeOpt{tfe.TeamOrganizationMemberships, tfe.TeamUsers},
		})
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tl.Items {
		fmt.Println(t.ID)
		for _, m := range t.Users {
			fmt.Println(m.Username)
		}
	}
}
