# Future GitHub Integration Changes

This note tracks the GitHub App work that should be added later after the current install/get/delete flow.

## 1. Installation Status

Add a `status` field to the GitHub installation record so the app can track:

- `active`
- `suspended`
- `uninstalled`

This will help the backend know whether the GitHub App is still usable.

## 2. Database Changes

Update the `github_installations` table to store the status value.

Planned changes:

- add a `status` column
- default new installations to `active`
- update the model in `domain/github.go`
- update insert/select SQL in the PostgreSQL repository

## 3. Suspend Handling

If GitHub marks the app as suspended later, do not delete the row immediately.

Instead:

- keep the installation record
- update status to `suspended`
- block actions that depend on an active installation

## 4. Uninstall Handling

If the user uninstalls the GitHub App, the backend should stop treating it as connected.

Later webhook support should:

- detect uninstall events from GitHub
- mark the installation as `uninstalled` or remove it
- keep the DB in sync automatically

## 5. Webhook Support

Webhook support will be added later so GitHub can notify the backend about installation changes.

Planned webhook events:

- installation created
- installation suspended
- installation unsuspended
- installation deleted

## 6. Handler Updates

After the status field is added, the SCM handler should return the stored status in the `GET` endpoint and use the status when deciding whether GitHub actions are allowed.

## 7. Suggested Follow-up Files

- `server/domain/github.go`
- `server/integrations/scm/github/repository/pgsql/pgsql_github.go`
- `server/integrations/scm/github/usecase/github_ucase.go`
- `server/integrations/scm/delivery/http/scm_handler.go`
- `server/migrations/*.sql`

