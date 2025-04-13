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
CLEAN="--clean"  # CLEAN é padrão agora
NO_CLEAN=false   # Nova flag para desativar o CLEAN

for arg in "$@"
do
  case $arg in
    --snapshot)
    SNAPSHOT="--snapshot"
    shift
    ;;
    --no-clean)
    NO_CLEAN=true
    shift
    ;;
    --clean)
    # Já é o padrão, mas mantemos para compatibilidade
    shift
    ;;
    *)
    # Unknown option
    ;;
  esac
done

# Desativa clean se solicitado explicitamente
if [ "$NO_CLEAN" = true ]; then
  CLEAN=""
fi

# Run GoReleaser
echo "Running GoReleaser..."
echo "Options: ${SNAPSHOT} ${CLEAN}"
goreleaser release $SNAPSHOT $CLEAN

# Unset the token for security
unset GITHUB_TOKEN
