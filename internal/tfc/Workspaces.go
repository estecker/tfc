package tfc

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-tfe"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
)

var PageSize int

func WorkspaceListCommand(TFCOrg string, Verbose bool, Search string, Tags string, ExcludeTags string, WildcardName string, Status string) error {
	PageSize = 10 //10 seems to be the sweet spot for speed
	wsl, err := WorkspaceList(TFCOrg, tfe.WorkspaceListOptions{
		Search:       Search,
		Tags:         Tags,
		ExcludeTags:  ExcludeTags,
		WildcardName: WildcardName,
		Include:      []tfe.WSIncludeOpt{tfe.WSCurrentRun},
		ListOptions:  tfe.ListOptions{PageSize: PageSize}})

	if err != nil {
		return fmt.Errorf("got a problem: %s", err)
	}
	for _, ws := range wsl {
		// If looking for a status, but not every workspace will actually have a current run
		if ws.CurrentRun != nil && tfe.RunStatus(Status) == ws.CurrentRun.Status {
			if Verbose {
				fmt.Printf("ID: %s, Name: %s, CRStatus: %s, Locked: %t\n", ws.ID, ws.Name, ws.CurrentRun.Status, ws.Locked)
			} else {
				fmt.Printf("%s\n", ws.ID) //Just want ws with status
			}
		} else if Status == "" {
			fmt.Printf("%s", ws.ID)
			if Verbose {
				if ws.CurrentRun != nil {
					fmt.Printf(": Name: %s, CRStatus: %s, Locked: %t", ws.Name, ws.CurrentRun.Status, ws.Locked)
				} else {
					fmt.Printf(":, Name: %s, Locked: %t", ws.Name, ws.Locked)
				}
			}
			fmt.Println("")
		}
	}
	return nil
}

func _(TFCOrg string) {
	_ = WorkspacePlanning(TFCOrg)
}
func _(workspaceID string) {
	_ = WorkspaceLock(workspaceID)
}
func _(workspaceID string) {
	_ = WorkspaceUnlock(workspaceID)
}

// concurrent workspace getter
func wsListWorker(client tfe.Client, in <-chan tfe.WorkspaceListOptions, out chan<- tfe.WorkspaceList, w8 *sync.WaitGroup, org string, errChan chan<- error) {
	defer w8.Done()
	for listOptions := range in {
		wsl, err := client.Workspaces.List(context.Background(), org, &listOptions)
		out <- *wsl
		if err != nil {
			errChan <- err
		}
	}
}

// WorkspaceList Get first page, then determine if more pages are needed
func WorkspaceList(org string, options tfe.WorkspaceListOptions) ([]tfe.Workspace, error) {
	var wsR []tfe.Workspace
	client, err := Client()
	ctx := context.Background()
	if err != nil {
		return nil, fmt.Errorf("could not create a client: %s", err)
	}
	wsl, err := client.Workspaces.List(ctx, org, &options)
	if err != nil {
		log.Fatalf("Could not list workspaces")
	}
	for _, ws := range wsl.Items {
		wsR = append(wsR, *ws)
	}
	if wsl.CurrentPage < wsl.TotalPages { //pagination time
		var w8 sync.WaitGroup //Wait group / channel creation
		w8.Add(wsl.TotalPages)
		inputCh := make(chan tfe.WorkspaceListOptions, 3) //Can work on 3 at once
		outputCh := make(chan tfe.WorkspaceList, wsl.TotalPages)
		errorCh := make(chan error)
		for k := 0; k < wsl.TotalPages; k++ { // Worker creation
			go wsListWorker(*client, inputCh, outputCh, &w8, org, errorCh)
		}
		// Already got page 1 above so start at page 2
		for p := 2; p <= wsl.TotalPages; p++ { // send tasks to the input channel
			options.ListOptions = tfe.ListOptions{PageNumber: p, PageSize: PageSize}
			inputCh <- options
		}
		close(inputCh)  // we say that we will not send any new data on the input channel
		w8.Wait()       // wait for all tasks to be completed
		close(outputCh) // when all treatment is finished we close the output channel
		close(errorCh)

		err := <-errorCh
		if err != nil {
			return nil, err
		}

		for out := range outputCh { //Collect the results
			for _, ws := range out.Items {
				wsR = append(wsR, *ws)
			}
		}
		return wsR, err
	} else {
		return wsR, err //no pagination
	}
	return nil, err
}

// WorkspacePlanning A workspace that's only planning, not ever applying
func WorkspacePlanning(org string) error {
	allWS, err := WorkspaceList(org, tfe.WorkspaceListOptions{})
	if err != nil {
		return fmt.Errorf("got a problem: %s", err)
	}
	client, err := Client()
	if err != nil {
		return fmt.Errorf("could not create a client: %s", err)
	}
	g, ctx := errgroup.WithContext(context.Background())
	for _, ws := range allWS {
		ws := ws
		g.Go(func() error {
			runList, err := client.Runs.List(ctx,
				ws.ID, &tfe.RunListOptions{
					ListOptions: tfe.ListOptions{PageSize: 1},
					Include:     []tfe.RunIncludeOpt{tfe.RunWorkspace},
					Operation:   "plan_only",
					Status:      "pending,fetching,fetching_completed,pre_plan_running,pre_plan_completed,queuing,plan_queued,planning,cost_estimating,policy_checking,post_plan_running",
					Source:      "tfe-api,tfe-configuration-version",
				})
			if err != nil {
				return err
			} else {
				for _, r := range runList.Items {
					fmt.Printf("I found a planonly run: %s in workspace: %s message: %+v\n", r.ID, r.Workspace.Name, r.Message)
				}
				return nil
			}
		})
	}
	if err := g.Wait(); err != nil {
		fmt.Printf("Error: %v", err)
	}
	return nil
}

func WorkspaceLock(workspaceID string) error {
	client, err := Client()
	ctx := context.TODO()
	if err != nil {
		return fmt.Errorf("could not create a client: %s", err)
	}
	reason := "Lock out, tag out"
	ws, err := client.Workspaces.Lock(ctx, workspaceID, tfe.WorkspaceLockOptions{Reason: &reason})
	if err == nil && ws.Locked {
		fmt.Printf("Locked workspace %s:%s\n", ws.ID, ws.Name)
		return nil
	} else {
		return fmt.Errorf("Had problems locking workspace %s:%s\n", ws.ID, ws.Name)
	}
}

func WorkspaceUnlock(workspaceID string) error {
	client, err := Client()
	ctx := context.TODO()
	if err != nil {
		return fmt.Errorf("could not create a client: %s", err)
	}
	ws, err := client.Workspaces.Unlock(ctx, workspaceID)
	if err == nil && ws.Locked == false {
		fmt.Printf("Unlocked workspace %s:%s\n", ws.ID, ws.Name)
		return nil
	} else {
		return fmt.Errorf("Had problems unlocking workspace %s:%s\n", ws.ID, ws.Name)
	}
}
