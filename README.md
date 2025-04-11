# Scopy - Smart Copy

Scopy é uma ferramenta de linha de comando escrita em Go que permite copiar o conteúdo de arquivos com extensões específicas de forma inteligente, respeitando configurações de exclusão e formatos personalizados.

## Instalação

```bash
go install github.com/dakoctba/scopy@latest
```

## Uso

```bash
scopy [opções] extensão1 extensão2 ...
```

### Opções

- `--header-format FORMAT` - Define o formato do cabeçalho que precede cada arquivo
  - O formato deve aceitar `%s` como placeholder para o caminho do arquivo
  - Exemplo: `--header-format "/* %s */"`
  - Valor padrão: `// %s`

- `--exclude PATTERNS` - Define padrões para excluir arquivos/diretórios
  - Aceita uma lista separada por vírgulas
  - Exemplo: `--exclude "node_modules,build,tmp"`

- `--list-only` - Apenas lista os arquivos que seriam copiados, sem mostrar seu conteúdo

- `--max-size SIZE` - Define o tamanho máximo dos arquivos a serem incluídos
  - Deve aceitar sufixos como KB, MB, GB
  - Exemplo: `--max-size 1MB`

- `--strip-comments` - Remove comentários dos arquivos de código
  - Deve suportar diferentes formatos de comentários com base na extensão do arquivo

- `--help` - Exibe ajuda e informações de uso

- `--version` - Exibe a versão do aplicativo

### Exemplos

```bash
# Copia arquivos .go e .js
scopy go js

# Customiza o formato do cabeçalho
scopy --header-format "/* %s */" go

# Ignora diretórios vendor e dist
scopy --exclude "vendor,dist" go js

# Lista apenas arquivos .go sem mostrar conteúdo
scopy --list-only go

# Ignora arquivos .go maiores que 500KB
scopy --max-size 500KB go

# Remove comentários dos arquivos copiados
scopy --strip-comments go js
```

## Estatísticas

Ao final da execução, o Scopy exibe estatísticas sobre os arquivos processados:
- Número total de arquivos processados
- Número de arquivos por extensão
- Tamanho total em bytes
- Número total de linhas copiadas

## Códigos de retorno

- 0: Execução bem-sucedida
- 1: Erro de uso (argumentos inválidos)
- 2: Erro ao ler/processar arquivos

## Licença

MIT License - veja o arquivo [LICENSE](LICENSE) para mais detalhes.
