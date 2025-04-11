# Prompt para desenvolvimento do aplicativo "scopy" (Smart Copy)

## Descrição do projeto
Desenvolver um aplicativo de linha de comando chamado "scopy" (Smart Copy) em Go que percorra recursivamente os subdiretórios do diretório atual, copiando o conteúdo de arquivos com extensões específicas para uma única saída consolidada.

## Funcionalidades principais
1. Percorrer recursivamente todos os subdiretórios a partir do diretório atual
2. Filtrar arquivos com base nas extensões fornecidas como argumentos
3. Copiar o conteúdo dos arquivos filtrados para uma única saída (stdout)
4. Adicionar um cabeçalho com o caminho relativo antes do conteúdo de cada arquivo
5. Respeitar regras do .gitignore se estiver em um repositório Git
6. Gerar estatísticas sobre os arquivos processados ao final da execução

## Flags e opções
O aplicativo deve suportar as seguintes flags:

1. `--header-format FORMAT` - Define o formato do cabeçalho que precede cada arquivo
   - O formato deve aceitar `%s` como placeholder para o caminho do arquivo
   - Exemplo: `--header-format "/* %s */"`
   - Valor padrão: `// %s`

2. `--exclude PATTERNS` - Define padrões para excluir arquivos/diretórios
   - Aceita uma lista separada por vírgulas
   - Exemplo: `--exclude "node_modules,build,tmp"`

3. `--list-only` - Apenas lista os arquivos que seriam copiados, sem mostrar seu conteúdo

4. `--max-size SIZE` - Define o tamanho máximo dos arquivos a serem incluídos
   - Deve aceitar sufixos como KB, MB, GB
   - Exemplo: `--max-size 1MB`

5. `--strip-comments` - Remove comentários dos arquivos de código
   - Deve suportar diferentes formatos de comentários com base na extensão do arquivo
   - Exemplo: para .go: remove `//`, `/* */`; para .js: remove `//`, `/* */`, etc.

6. `--help` - Exibe ajuda e informações de uso

7. `--version` - Exibe a versão do aplicativo

## Comportamento de entrada/saída
- O aplicativo deve ler de stdin quando necessário e escrever para stdout
- Mensagens de erro devem ser escritas em stderr
- Deve seguir o princípio do Unix de "fazer uma coisa e fazer bem"
- O resultado deve ser facilmente redirecionável para um arquivo ou encadeado com outros comandos

## Integração com .gitignore
- Se o diretório atual for um repositório Git com um arquivo .gitignore
- Processar as regras do .gitignore e ignorar arquivos/diretórios que correspondam aos padrões

## Estatísticas de saída
O aplicativo deve automaticamente gerar estatísticas ao final da execução, incluindo:
- Número total de arquivos processados
- Número de arquivos por extensão
- Tamanho total em bytes
- Número total de linhas copiadas

## Códigos de retorno
- 0: Execução bem-sucedida
- 1: Erro de uso (argumentos inválidos)
- 2: Erro ao ler/processar arquivos

## Estrutura do código
- Separar em módulos/pacotes para melhor organização e testabilidade
- Utilizar interfaces Go para facilitar testes e extensões futuras
- Seguir convenções de nomeação e estruturação do Go

## Opções de uso
```
scopy [opções] extensão1 extensão2 ...

Exemplos:
  scopy go js                               # Copia arquivos .go e .js
  scopy --header-format "/* %s */" go       # Customiza o formato do cabeçalho
  scopy --exclude "vendor,dist" go js       # Ignora diretórios vendor e dist
  scopy --list-only go                      # Lista apenas arquivos .go sem mostrar conteúdo
  scopy --max-size 500KB go                 # Ignora arquivos .go maiores que 500KB
  scopy --strip-comments go js              # Remove comentários dos arquivos copiados
```

## Documentação
- Incluir comentários detalhados no código
- Fornecer exemplos de uso
- Criar uma página man ou documentação equivalente

## Testes
- Desenvolver testes unitários para cada módulo/função
- Testes de integração para comportamentos complexos (ex: processamento de regras .gitignore)
- Testes de borda para verificar casos limites (arquivos vazios, arquivos grandes, etc.)

## Validações importantes
- Validar existência do diretório atual
- Validar formatos e valores das flags
- Lidar adequadamente com erros de leitura de arquivo
- Tratar corretamente caminhos de arquivos em diferentes sistemas operacionais

## Licença
- MIT License ou similar (Open Source)

## Observações finais
O aplicativo deve seguir as boas práticas de ferramentas Unix-like, enfatizando simplicidade, modularidade, e a capacidade de interagir bem com outras ferramentas através de pipes e redirecionamentos.
