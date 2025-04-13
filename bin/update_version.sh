#!/bin/bash

# Cores para saída
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Lê a versão atual do arquivo VERSION
CURRENT_VERSION=$(cat VERSION | tr -d '[:space:]')
echo -e "${BLUE}Versão atual:${NC} ${GREEN}$CURRENT_VERSION${NC}"

# Sugere a próxima versão incrementando o patch
IFS='.' read -r major minor patch <<< "$CURRENT_VERSION"
NEXT_PATCH=$((patch + 1))
NEXT_MINOR=$((minor + 1))
NEXT_MAJOR=$((major + 1))

echo -e "\n${YELLOW}Sugestões para próxima versão:${NC}"
echo -e "  1. Patch (correção de bugs): ${GREEN}$major.$minor.$NEXT_PATCH${NC}"
echo -e "  2. Minor (novos recursos): ${GREEN}$major.$NEXT_MINOR.0${NC}"
echo -e "  3. Major (mudanças incompatíveis): ${GREEN}$NEXT_MAJOR.0.0${NC}"
echo -e "  4. Personalizada (digite você mesmo)"

# Perguntar qual tipo de versão o usuário quer
read -p $'\e[1;33mEscolha uma opção (1-4): \e[0m' VERSION_OPTION

case $VERSION_OPTION in
  1)
    NEW_VERSION="$major.$minor.$NEXT_PATCH"
    ;;
  2)
    NEW_VERSION="$major.$NEXT_MINOR.0"
    ;;
  3)
    NEW_VERSION="$NEXT_MAJOR.0.0"
    ;;
  4)
    read -p $'\e[1;33mDigite a nova versão: \e[0m' NEW_VERSION
    ;;
  *)
    echo -e "${RED}Opção inválida!${NC}"
    exit 1
    ;;
esac

# Confirmar a nova versão
echo -e "\n${BLUE}Resumo da operação:${NC}"
echo -e "  De: ${GREEN}$CURRENT_VERSION${NC}"
echo -e "  Para: ${GREEN}$NEW_VERSION${NC}"
echo -e "\nIsto irá:"
echo -e "  1. Atualizar o arquivo VERSION"
echo -e "  2. Criar uma tag git v$NEW_VERSION"
echo -e "  3. Fazer push da tag para o repositório remoto"
echo -e "  4. Gerar uma release usando o GoReleaser"

read -p $'\e[1;33mConfirma esta operação? (S/n): \e[0m' CONFIRM
CONFIRM=${CONFIRM:-S} # Valor padrão é S

if [[ $CONFIRM =~ ^[Ss]$ ]]; then
  # Atualizar o arquivo VERSION
  echo "$NEW_VERSION" > VERSION
  echo -e "\n${GREEN}✓ Arquivo VERSION atualizado para $NEW_VERSION${NC}"

  # Commit da alteração do arquivo VERSION
  git add VERSION
  git commit -m "Bump version to $NEW_VERSION"
  echo -e "${GREEN}✓ Alterações commitadas${NC}"

  # Executar o script de release
  echo -e "\n${BLUE}Executando processo de release...${NC}"

  # Chamar o script de release passando os argumentos
  # O script release.sh vai cuidar de criar a tag e fazer o push
  $(dirname "$0")/release.sh

  echo -e "\n${GREEN}Processo de versão e release concluído com sucesso!${NC}"
else
  echo -e "\n${YELLOW}Operação cancelada pelo usuário.${NC}"
  exit 0
fi
