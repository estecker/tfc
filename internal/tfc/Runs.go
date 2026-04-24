package tfc

import (
	"bufio"
	"context"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"log"
	"os"
	"strings"
	"time"
)

func RunsCmd(Workspaces []string, Force bool, viewCurrentRun bool) {
	for _, ws := range Workspaces {
		workspaceRun(ws, Force, viewCurrentRun)
	}
}

func workspaceRun(wsid string, force bool, viewCurrentRun bool) {
	client, err := Client()
	ctx := context.Background()
	if err != nil {
		log.Fatalf("Issues creating a client\n")
	}
	if viewCurrentRun {
		ws, err := client.Workspaces.ReadByIDWithOptions(ctx, wsid, &tfe.WorkspaceReadOptions{Include: []tfe.WSIncludeOpt{tfe.WSCurrentRun}})
		if err != nil {
			log.Fatalf("ERROR ReadByIDWithOptions\n")
		}
		plan, err := client.Plans.Read(ctx, ws.CurrentRun.Plan.ID)
		if err != nil {
			log.Fatalf("Failed to retrieve plan %s\n%s", ws.CurrentRun.Plan.ID, err)
		}
		if plan.Status == tfe.PlanErrored {
			fmt.Printf("The plan %q errored, no sense on trying again?\n", plan.ID)
			renderLog(client.Plans.Logs(ctx, ws.CurrentRun.Plan.ID))
		} else {
			fmt.Printf("Found a workspace that named: %s\nstatus: %s\nmessage: %s\nrunID: %s\napplyID: %s\n\n", ws.Name, ws.CurrentRun.Status, ws.CurrentRun.Message, ws.CurrentRun.ID, ws.CurrentRun.Apply.ID)
			if ws.CurrentRun.Status == tfe.RunApplied || ws.CurrentRun.Status == tfe.RunErrored {
				renderLog(client.Applies.Logs(ctx, ws.CurrentRun.Apply.ID))
				fmt.Printf("HasChanges: %t, ResourceChanges: %d, ResourceAdditions: %d, ResourceDestructions: %d\n", plan.HasChanges, plan.ResourceChanges, plan.ResourceAdditions, plan.ResourceDestructions)
			} else {
				fmt.Printf("Current run probably does not have any logs to show\n")
			}
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Start a new run [y/N]: ")

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" {
			os.Exit(0)
		}
	}

	message := time.Now().Format(time.RFC3339) + " Don't communicate by sharing memory, share memory by communicating."
	run, err := client.Runs.Create(ctx, tfe.RunCreateOptions{Workspace: &tfe.Workspace{ID: wsid}, Message: &message})
	if err != nil {
		log.Fatalf("ERROR creating run")
	}
	ws, err := client.Workspaces.ReadByID(ctx, wsid) //let's find org from workspaceID rather than require org cmd arg
	baseURL := client.BaseURL()
	fmt.Printf("%s://%s/app/%s/workspaces/%s/runs/%s\n", baseURL.Scheme, baseURL.Host, ws.Organization.Name, ws.Name, run.ID)
	if force == true {
		ws, _ := client.Workspaces.ReadByIDWithOptions(ctx, wsid, &tfe.WorkspaceReadOptions{Include: []tfe.WSIncludeOpt{tfe.WSCurrentRun}})
		if ws.CurrentRun.ID != run.ID && run.Status == tfe.RunPending {
			fmt.Printf("Will ForceExecute %s because it appears that %s is stuck in planned state.\n", run.ID, ws.CurrentRun.ID)
			_ = client.Runs.ForceExecute(ctx, run.ID)
		} else {
			fmt.Printf("It looks like the new run %s does not need to be forced.\n", run.ID)
		}
	}
}

func RunsDiscardCmd(Workspace string) {
	workspaceDiscard(Workspace)
}

// Just discard but not start anything new
func workspaceDiscard(wsid string) {
	client, err := Client()
	ctx := context.TODO()
	if err != nil {
		log.Fatalf("Issues creating a client")
	}
	//ws, err := client.Workspaces.ReadByIDWithOptions(ctx, wsid, &tfe.WorkspaceReadOptions{Include: []tfe.WSIncludeOpt{tfe.WSCurrentRun}})
	//if err != nil {
	//	fmt.Printf("ERROR ReadByIDWithOptions\n")
	//}
	//if ws.CurrentRun.Status == tfe.RunPlanned || ws.CurrentRun.Status == tfe.RunPolicyChecked {
	//	fmt.Printf("Will discard run %s in %s.\n", ws.CurrentRun.ID, ws.Name)
	//	message := "The discard action can be used to skip any remaining work on runs that are paused waiting for confirmation or priority. Miłego dnia."
	//	_ = client.Runs.Discard(ctx, ws.CurrentRun.ID, tfe.RunDiscardOptions{Comment: &message})
	//} else {
	//	fmt.Printf("It does not look like there's anything to discard. Current run status is: %s:\n", ws.CurrentRun.Status)
	//}

	// This is a bit more aggressive, it will cancel all pending runs
	message := "The discard action can be used to skip any remaining work on runs that are paused waiting for confirmation or priority. Miłego dnia."
	runs, err := client.Runs.List(ctx, wsid, &tfe.RunListOptions{Status: "pending"})
	for _, run := range runs.Items {
		fmt.Printf("Will cancel run %s\n", run.ID)
		//client.Runs.Discard(ctx, run.ID, tfe.RunDiscardOptions{Comment: &message})
		client.Runs.Cancel(ctx, run.ID, tfe.RunCancelOptions{Comment: &message})
	}
}
func RunsLogCmd(RunID string) {
	runLog(RunID)
}

func runLog(runid string) {
	ctx := context.TODO()
	client, err := Client()
	if err != nil {
		log.Fatalf("Issues creating a client")
	}
	run, err := client.Runs.ReadWithOptions(ctx, runid, &tfe.RunReadOptions{Include: []tfe.RunIncludeOpt{tfe.RunPlan}})
	if err != nil {
		log.Fatalf("Could not find a plan for runid: %s\n", runid)
	}
	renderLog(client.Plans.Logs(ctx, run.Plan.ID))
	if err != nil {
		log.Fatalf("Failed to retrieve plan %s\n", runid)
	}

}

func RunCancelCmd(RunID string) {
	runCancel(RunID)
}
func runCancel(runid string) {
	ctx := context.TODO()
	client, err := Client()
	if err != nil {
		log.Fatalf("Issues creating a client")
	}
	comment := time.Now().Format(time.RFC3339) + " Time to cancel"
	err = client.Runs.Cancel(ctx, runid, tfe.RunCancelOptions{Comment: &comment})
	if err != nil {
		log.Fatalf("Failed to cancel %s\n", runid)
	}
}
