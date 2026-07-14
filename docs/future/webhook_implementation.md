# Webhook Implementation Review Plan

Updated: 15 July 2026

This note captures what the webhook layer still needs in the current `server/` architecture.

## Current Shape

- The GitHub webhook parser exists and already verifies signatures when a secret is configured.
- The webhook route still needs to be wired into the API server.

## What The Webhook Still Needs

- A dedicated HTTP route wired into the server.
- A handler that accepts GitHub installation lifecycle events.
- Logic for `installation_suspend`, `installation_unsuspend`, and `installation.deleted`.
- A clear rule for whether reinstall events update the existing row or create a new one.
- A decision on how much of the repo-selection flow should be added before webhook work is considered complete.

## What The Webhook Should Do

- Verify GitHub webhook authenticity before parsing the payload.
- Keep installation status in sync with GitHub's lifecycle events.
- Mark suspended installations as `suspended`.
- Mark removed installations as `uninstalled`.
- Leave local deployment history and application data intact.
- Block GitHub API-dependent actions when installation status is not `active`.

## Recommended Boundaries

- Installation means access.
- Repository listing means selection.
- Webhook means state synchronization.

Keeping those boundaries separate will make the deployment flow much easier to extend.

## Suggested Next Steps

1. Wire the webhook route into `server/app/main.go`.
2. Add a GitHub webhook handler in the SCM area.
3. Connect lifecycle events to the installation status update path.
4. Add repository listing for installed accounts.
5. Add repository selection storage before deployment expands further.
