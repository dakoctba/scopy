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

# Get version from VERSION file
VERSION=$(cat VERSION | tr -d '[:space:]')
TAG="v$VERSION"
echo "Using version $VERSION (tag $TAG) from VERSION file"

# Parse arguments
SNAPSHOT=""
CLEAN="--clean"  # CLEAN é padrão agora
NO_CLEAN=false   # Nova flag para desativar o CLEAN
SKIP_TAG=false   # Opção para pular criação/verificação de tag

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
    --skip-tag)
    SKIP_TAG=true
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

# Cria tag se ela não existir e não estamos no modo snapshot e não foi pedido para pular
if [ -z "$SNAPSHOT" ] && [ "$SKIP_TAG" != true ]; then
  if ! git rev-parse "$TAG" >/dev/null 2>&1; then
    echo "Creating git tag $TAG..."
    git tag -a "$TAG" -m "Version $VERSION"
    echo "Pushing tag to remote..."
    git push origin "$TAG"
  else
    echo "Tag $TAG already exists, using existing tag"
  fi
fi

# Run GoReleaser
echo "Running GoReleaser..."
echo "Options: ${SNAPSHOT} ${CLEAN}"
goreleaser release $SNAPSHOT $CLEAN

# Unset the token for security
unset GITHUB_TOKEN
