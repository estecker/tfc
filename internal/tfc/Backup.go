package tfc

import (
    "fmt"
    "github.com/hashicorp/go-tfe"
    "log"
    "path/filepath"
)

func BackupStatesCmd(Org string, Folder string) error {
    wsl, err := WorkspaceList(Org, tfe.WorkspaceListOptions{
        ListOptions: tfe.ListOptions{PageSize: PageSize},
        Include:     []tfe.WSIncludeOpt{tfe.WSCurrentStateVer},
    })
    if err != nil {
        log.Fatal("got a problem: %s", err)
    }
    total := len(wsl)
    index := total
    for _, ws := range wsl {
        if ws.CurrentStateVersion != nil {
            fullPath := filepath.Join(Folder, ws.Name+".json")
            err := downloadStateFile(fullPath, ws.CurrentStateVersion.DownloadURL)
            if err != nil {
                return fmt.Errorf(" Could not download the state for %s: %s\n", ws.Name, err)
            }
        } else {
            fmt.Printf("  No state for %s\n", ws.Name)
        }
        fmt.Printf("\033[2K\r%d", index) // Clear the current line, return the cursor to the beginning, and print the index
        index--
    }
    fmt.Printf("\n")
    return nil
}

func BackupVariablesCmd(Org string, Folder string) error {
    wsl, err := WorkspaceList(Org, tfe.WorkspaceListOptions{})
    if err != nil {
        return fmt.Errorf("BackupVariablesCmd got a problem: %s", err)
    }
    total := len(wsl)
    index := total
    for _, ws := range wsl {
        allVars, err := VariableList(ws.ID, tfe.VariableListOptions{ListOptions: tfe.ListOptions{PageSize: 100}}) //TODO pagination
        if err != nil {
            return fmt.Errorf("Could not get variables from %s\n", ws.Name)
        }
        fullPath := filepath.Join(Folder, ws.Name+".json")
        for _, v := range allVars.Items {
            v.Workspace = nil //Don't need this info in the backup file
        }
        err = downloadStruct(fullPath, allVars.Items)
        if err != nil {
            return fmt.Errorf("error problem writing the file: %s", err)
        }
        fmt.Printf("\033[2K\r%d", index) // Clear the current line, return the cursor to the beginning, and print the index
        index--
    }
    fmt.Printf("\n")
    return nil
}

func BackupWorkspacesCmd(Org string, Folder string) error {
    wsl, err := WorkspaceList(Org, tfe.WorkspaceListOptions{})
    if err != nil {
        return fmt.Errorf("could not list workspaces: %w", err)
    }
    total := len(wsl)
    index := total
    log.Println("Backing up workspaces")
    log.Println("Total workspaces: ", total)
    for _, ws := range wsl {
        fullPath := filepath.Join(Folder, ws.Name+".json")
        err := downloadStruct(fullPath, &ws)
        if err != nil {
            return fmt.Errorf("could not write the file: %w", err)
        }
        fmt.Printf("\033[2K\r%d", index) // Clear the current line, return the cursor to the beginning, and print the index
        index--
    }
    fmt.Printf("\n")
    return nil
}

