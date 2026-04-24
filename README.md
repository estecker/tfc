# Goals
* Re-run all failed workspaces
  * `for ws in $(tfc workspaces --org <yourorg> --status errored );do tfc runs --workspace-ids ${ws};done`
* Clear all pending plan/applies, then trigger a fresh plan/apply
  * `for ws in $(tfc workspaces --org <yourorg> --status pending );do tfc runs --workspace-ids ${ws} --force;done`

* Discard all pending plans/applies
* `for ws in $(tfc workspaces --org <yourorg> --status policy_checked --exclude-tags "slice:legacy");do tfc runs discard --workspace-id ${ws};done`
* 
* Apply all workspaces in RunPolicyChecked state
  * `for ws in $(tfc workspaces --org <yourorg> --status policy_checked --exclude-tags "slice:legacy");do tfc runs --workspace-ids ${ws} --force;done`

* List planonly runs currently running, detect large PR
  * `tfc workspaces --org <yourorg> planning`

* List workspaces actively working on something

* List workspaces blocked on something


* Cancel plan in large PR, especially ones far behind in commits
  * `tfc runs cancel --run-id run-c6vzYfGctBNxEVoh`

* - [ ]  List locked workspaces, detect large operation (Merged PR)
  * `tfc workspaces --org <yourorg>  --verbose | grep "Locked: true"`
  
* - [ ] Lock a workspace
 * `tfc workspaces lock --id ws-wC89HymkmQ4Botg4`
* - [ ] Unlock a workspace
 * `tfc workspaces unlock --id ws-wC89HymkmQ4Botg4`
 
* - [ ] Drift detection TODO no API support
TODO refresh?


* - [ ] Backups
  `tfc --org <yourorg> backup states    --folder "${WORKDIR}/states/"`
  `tfc --org <yourorg> backup variables --folder "${WORKDIR}/variables/"`
  `tfc --org <yourorg> backup workspace --folder "${WORKDIR}/workspaces/"`


* - [ ]  Non use cases
  * - [ ]  Settings
  * - [ ]  Or other non-operational changes

