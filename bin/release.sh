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
SKIP_TAG=false   # Opção para pular criação/verificação de tag
VERSION_ARG=""

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
    --version=*)
    VERSION_ARG="${arg#*=}"
    shift
    ;;
    *)
    # Unknown option
    ;;
  esac
done

# Determina a versão a ser usada
if [ -n "$VERSION_ARG" ]; then
  # Usa a versão fornecida como argumento
  VERSION="$VERSION_ARG"
else
  # Obtém a versão da última tag git
  VERSION=$(git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')
  if [ -z "$VERSION" ]; then
    echo "Error: No version specified and no git tags found"
    exit 1
  fi
fi

TAG="v$VERSION"
echo "Using version $VERSION (tag $TAG)"

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
