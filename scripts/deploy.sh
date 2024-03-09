#!/bin/bash

# Install flyctl if isn't present
if ! command -v flyctl &> /dev/null; then
    curl -L https://fly.io/install.sh | sh
    export FLYCTL_INSTALL="/home/runner/.fly"
    export PATH="$FLYCTL_INSTALL/bin:$PATH"
fi

# Deploy
flyctl deploy --remote-only \
  --build-secret UNCONDITIONAL_API_SOURCE_CLIENT_KEY="$UNCONDITIONAL_API_SOURCE_CLIENT_KEY" \
  --build-secret UNCONDITIONAL_API_FEED_REPO_KEY="$UNCONDITIONAL_API_FEED_REPO_KEY" \
  --build-secret UNCONDITIONAL_API_FEED_REPO_HOST="$UNCONDITIONAL_API_FEED_REPO_HOST" \
  --build-arg UNCONDITIONAL_API_FEED_REPO_INDEX="$UNCONDITIONAL_API_FEED_REPO_INDEX" \
  --build-arg UNCONDITIONAL_API_BUILD_COMMIT_VERSION="$UNCONDITIONAL_API_BUILD_COMMIT_VERSION" \
  --build-arg UNCONDITIONAL_API_BUILD_RELEASE_VERSION="$UNCONDITIONAL_API_BUILD_RELEASE_VERSION"
