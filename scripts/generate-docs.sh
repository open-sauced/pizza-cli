#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

# Configure git to use the OpenSauced bot account
git config user.name 'open-sauced[bot]'
git config user.email '63161813+open-sauced[bot]@users.noreply.github.com'

# Semantic release made changes, so pull the latest changes from the current branch
git pull origin "$GITHUB_REF"

# Generate documentation
just gen-docs

# Get the author of the last non-merge commit
LAST_COMMIT_AUTHOR=$(git log -1 --no-merges --pretty=format:'%an <%ae>')

# Commit with co-authorship and push changes
git add docs/
git commit -m "chore: automated docs generation for release

Co-authored-by: $LAST_COMMIT_AUTHOR" || echo "No changes to commit"
git push origin HEAD:"$GITHUB_REF"
