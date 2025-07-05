#!/bin/bash

# Check if new repo name is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <new-repo-name>"
    exit 1
fi

NEW_REPO_NAME=$1
OLD_REPO_NAME="go-template"

# Find all .go files and replace the old repo name with the new one
find . -type f -name "*.go" -exec sed -i "s|$OLD_REPO_NAME|$NEW_REPO_NAME|g" {} +

# Update go.mod file
sed -i "s|module $OLD_REPO_NAME|module $NEW_REPO_NAME|g" go.mod

echo "Repository name has been updated from '$OLD_REPO_NAME' to '$NEW_REPO_NAME'"
