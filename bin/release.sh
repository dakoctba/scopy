#!/bin/bash

# Check if .env file exists
if [ -f .env ]; then
  echo "Loading environment variables from .env file"
  export $(grep -v '^#' .env | xargs)
else
  echo "Warning: .env file not found"
fi

# Check if GITHUB_TOKEN is set
if [ -z "$GITHUB_TOKEN" ]; then
  echo "Error: GITHUB_TOKEN is not set"
  echo "Please set it in your .env file or export it manually"
  exit 1
fi

# Parse arguments
SNAPSHOT=""
CLEAN=""

for arg in "$@"
do
  case $arg in
    --snapshot)
    SNAPSHOT="--snapshot"
    shift
    ;;
    --clean)
    CLEAN="--clean"
    shift
    ;;
    *)
    # Unknown option
    ;;
  esac
done

# Run GoReleaser
echo "Running GoReleaser..."
goreleaser release $SNAPSHOT $CLEAN

# Unset the token for security
unset GITHUB_TOKEN
