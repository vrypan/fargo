#!/bin/bash

# Calculate ne tag version
git fetch --tags
LATEST_TAG=$(git tag --sort=-v:refname | head -n 1)
LATEST_TAG=${LATEST_TAG:-"v0.0.0"}

VERSION=$(echo $LATEST_TAG | sed 's/v//')
IFS='.' read -r -a PARTS <<< "$VERSION"

# Default to patch increment
case "$1" in
  major)
    PARTS[0]=$((PARTS[0] + 1))
    PARTS[1]=0
    PARTS[2]=0
    ;;
  minor)
    PARTS[1]=$((PARTS[1] + 1))
    PARTS[2]=0
    ;;
  patch|""|*)
    PARTS[2]=$((PARTS[2] + 1))
    ;;
esac

NEW_TAG="v${PARTS[0]}.${PARTS[1]}.${PARTS[2]}"
echo "Creating new tag: $NEW_TAG"

git tag -a "$NEW_TAG" -F tmp/release_notes.md
git push origin "$NEW_TAG"

echo "Tag $NEW_TAG pushed successfully!"
