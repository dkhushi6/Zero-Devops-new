# Future Issues And Priorities

This file lists known follow-up issues after the current auth and GitHub App installation APIs were manually tested.

Current decision:

- Do not change the current delete flow right now.
- Do not add installation `status` right now.
- Treat status and webhook handling as later work.
- Move next into repository access and deployment flow.

## Important

- Build the next product flow after GitHub App installation: list the repositories available to the installed GitHub App, let the user choose a repository, then prepare that repository for deployment.
- Use the stored GitHub `installation_id` to create a GitHub App installation access token before reading repository data.
- Add a repository listing usecase and API for the authenticated user.
- Confirm GitHub App permissions include the repository access needed for reading repository metadata and code.
- Verify the deployed environment end to end on the real domain after auth and installation are connected to repository selection.

## Strict

- Keep the GitHub OAuth login flow separate from the GitHub App installation flow.
- Keep using middleware context for `user_id`; do not store request-specific user state on handlers.
- Keep protected SCM APIs behind the access-token middleware.
- Before cloning or deploying a repository, verify that the selected repository belongs to the authenticated user's GitHub App installation.
- Do not assume installation means access to all repositories; use the installation token and GitHub's allowed repository list.

## Optional

- Rename routes later for cleaner API shape, for example `/integrations/github/installation` instead of `/integration/scm/github/`.
- Add `INSERT ... ON CONFLICT ... DO UPDATE` for GitHub installation storage if reinstalling should update the same user row.
- Improve error-to-status-code mapping in handlers so auth and SCM errors return more precise `400`, `401`, `404`, and `409` responses.
- Hash refresh tokens before production instead of storing raw refresh tokens.
- Replace or refresh older documentation sections in `revise.md` that describe now-completed compile issues.

## Mild

- Normalize spelling and naming in older docs where `installation` was previously misspelled.
- Consider renaming `GET /auth/user/me` to `GET /auth/me` later for a smaller public API.
- Clean up older historical notes in `revise.md` once the project has a dedicated changelog.
- Improve response messages for GitHub App install/delete to use consistent capitalization.

## Later: Status And Webhooks

These are deliberately not part of the immediate next implementation.

- Add `status` to `github_installations` later if the app needs to preserve `active`, `suspended`, and `uninstalled` states.
- Add GitHub webhook handling later for installation events.
- Listen for installation suspend, unsuspend, and uninstall events.
- Decide later whether uninstall should delete the row or mark it as `uninstalled`.
- Update repository and handler responses later to include installation status.

## Recommended Next Implementation

The next implementation should be repository access for deployments.

Planned flow:

1. User logs in with GitHub OAuth.
2. User installs the GitHub App.
3. Backend stores the GitHub App `installation_id`.
4. Backend creates an installation access token using that `installation_id`.
5. Backend lists repositories available to that installation.
6. User selects one repository.
7. Backend fetches repository metadata and prepares the deployment flow.

Suggested first API:

```text
GET /integration/scm/github/repos
```

Expected behavior:

- Requires `access_token` cookie.
- Reads `user_id` from middleware context.
- Loads the stored GitHub installation for that user.
- Creates a GitHub App installation token.
- Calls GitHub to list repositories accessible to the installation.
- Returns the repository list to the frontend.

