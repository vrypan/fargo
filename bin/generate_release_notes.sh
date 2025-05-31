#!/bin/bash

# Generate Release Notes
PREV_TAG=$(git describe --tags --abbrev=0)
echo "Generating release notes from $PREV_TAG to HEAD..."
echo "## Changelog" > tmp/release_notes.md
git log --pretty=format:"- %s" $PREV_TAG..HEAD >> tmp/release_notes.md
cat tmp/release_notes.md
