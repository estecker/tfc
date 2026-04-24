package tfc

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"log"
)

func OrganizationCmd(org string) {
	var r []*tfe.OrganizationMembership
	r, err := orgMembership(org)
	if err != nil {
		log.Fatalf("Issues listing organization memberships\n")
	}
	for _, om := range r {
		fmt.Printf("%s:%s", om.User.ID, om.Email)
		for _, t := range om.Teams {
			fmt.Printf(",%s", t.Name)
		}
		fmt.Println()
	}

}

func orgMembership(org string) ([]*tfe.OrganizationMembership, error) {
	client, err := Client()
	ctx := context.Background()
	if err != nil {
		log.Fatalf("Issues creating a client\n")
	}
	pageNumber := 1
	var r []*tfe.OrganizationMembership
	for {
		oml, err := client.OrganizationMemberships.List(ctx,
			org,
			&tfe.OrganizationMembershipListOptions{ListOptions: tfe.ListOptions{PageNumber: pageNumber},
				Include: []tfe.OrgMembershipIncludeOpt{tfe.OrgMembershipUser, tfe.OrgMembershipTeam}})
		if err != nil {
			log.Fatalf("Issues listing organization memberships\n")
			return nil, err
		}
		for _, m := range oml.Items {
			r = append(r, m)
		}
		// Check if we've reached the last page
		if oml.Pagination.NextPage == 0 {
			break
		}
		// Move to the next page
		pageNumber = oml.Pagination.NextPage
	}
	return r, nil
}
