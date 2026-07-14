# Future Issues And Priorities

Updated: 15 July 2026

This file tracks what is still left after the current auth, installation, deployment, and worker flows are in place.

## Highest Priority

- Add a repository-listing API so the frontend can show repositories available to the installed GitHub App.
- Keep repository selection separate from installation storage.
- Use the stored GitHub `installation_id` to create an installation access token before reading repo data.
- Verify selected repositories actually belong to the authenticated user's installation.
- Wire the webhook endpoint so installation lifecycle events can update local state.

## Important

- Keep the current auth flow separate from the GitHub App installation flow.
- Add the missing webhook lifecycle handlers for `installation_suspend`, `installation_unsuspend`, and `installation.deleted`.
- Decide whether installation rows should be updated or soft-deleted on uninstall based on product needs.

## Useful Follow-Ups

- Rename routes later for a cleaner public API, for example `/integrations/github/installation` instead of `/integration/scm/github/`.
- Improve error-to-status-code mapping so auth, SCM, and deployment errors return more precise `400`, `401`, `404`, and `409` responses.
- Hash refresh tokens before production instead of storing raw refresh tokens.
- Tighten worker retry and recovery behavior if you want more reliable builds.
- Revisit older notes in `revise.md` as the new docs become the source of truth.

## Later

- Add webhook-driven installation status updates.
- Keep local deployment history, build logs, and stored app data even if GitHub access becomes inactive.
- Block GitHub API-dependent actions when installation status is not `active`.
- Treat reinstall events as reconnection events and refresh the existing row instead of blindly creating duplicates.
- Add deployment-time repository listing, branch selection, and sync checks after the repo picker exists.
