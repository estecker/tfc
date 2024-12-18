package tfc

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"log"
)

func VariableList(workspaceID string, options tfe.VariableListOptions) (tfe.VariableList, error) {
	client, err := Client()
	ctx := context.TODO()
	if err != nil {
		return tfe.VariableList{}, fmt.Errorf("could not create a client: %s", err)
	}
	varLR := tfe.VariableList{}
	for {
		varList, err := client.Variables.List(ctx, workspaceID, &options)
		if err != nil {
			log.Fatalf("Could not list variables")
		}
		for _, ws := range varList.Items {
			varLR.Items = append(varLR.Items, ws)
		}
		if varList.CurrentPage < varList.TotalPages {
			options.ListOptions = tfe.ListOptions{PageNumber: varList.NextPage}
		} else {
			return varLR, err
		}
	}
	return varLR, err
}
