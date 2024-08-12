# Ganho de capital

## Decisoes tecnicas e arquiteturais
- Utilizei a linguagem Go por ser uma linguagem que estou aprendendo recentemente e por estar gostando bastante de suas features, além de aproveitar o challenge para treinar.

  Dentre as features que achei relevantes para o desafio, destacam-se:

    1. Maior controle sobre alocacao de memoria, quando comparado com Ruby ou Javascript (as outras linguagens que domino)
    2. Facilidade em trabalhar com concorrencia com goroutines (caso seja interessante em alguma iteracao futura)
    3. Tipagem estatica, o que facilita muito no quesito validacao de inputs
 - Sobre a arquitetura da aplicacao, utilizei um padrao conhecido na comunidade Go, onde existe uma cada de binarios executaveis (/cmd), resposavel somente por orquestrar as outras camadas internas da aplicacao e direcionar input e output, alem da organizacao atraves de pacotes internals, uma feature do próprio Go que torna os pacotes contidos dentro dele privados daquela aplicacao em especifico, servindo assim de casa para toda a regra de negocio da aplicacao.
 - No quesito Docker, optei por uma imagem mais simples e somente para desenvolvimento, nao estando apta para ser usada em producao por alguns motivos:

    1. Imagem sem muilti-stage builds, o que torna a sua build mais lenta, alem de gerar uma imagem de tamanho maior desnecessaria
    2. Uma imagem final baseada em `scrach` ou `alpine` acaba nao sendo eficaz para rodar testes ou ate mesmo rodar o `sh` como command, o que tornaria o processo de compilacao e testes mais complexo
 
 ## Justificativa de libs 3rd party
 Neste projeto utilizei somente 1 lib 3rd party (github.com/shopspring/decimal).
 Julguei condizente com o challenge de utilizar esta biblioteca pelos seguintes motivos:
  
  1. Em um sistema de producao real, lidar com aritmetica de decimais com precisao arbitraria e algo recorrente em aplicacoes financeiras, logo faz muito sentido abracarmos o open source e utilizarmos o conhecimento compartilhado de anos de battle test das libs mais utilizadas (como o caso da lib em questao)
  2. Acreditei estar fora do escopo deste challenge implementar uma lib que lide com maestria todos os casos de borda para se trabalhar com aritmetica de numeros flutuantes, seja qual for a tecnica, entao para contemplar esse quesito e facilitar o desenvolvimento, optei por utilizar a biblioteca

## Como compilar e executar o projeto
Esta aplicacao conta com um `Dockerfile` para facilitar este processo.
Seguindo os seguintes passos, em um terminal, conseguimos testar a aplicacao (partindo do pressuposta que ja tenha um Docker engine rodando em sua maquina host):
```bash
# Dentro da raiz do projeto
docker build . -t stocks-tax-calculator
docker run -it stocks-tax-calculator sh

# com o sh rodando
/stocks-tax-calculator < input.txt

# a seguinte saida sera printada no terminal
[{"tax":0.00},{"tax":0.00},{"tax":0.00}]
[{"tax":0.00},{"tax":10000.00},{"tax":0.00}]
```
Junto do repositorio, em sua raiz, existe um arquivo de teste chamado `input.txt` com dados para teste, caso queria testar mais casos diferentes, altere o `input.txt` e repita o processo.

A saida esperada do arquivo de exemplo se encontra no `results.txt`, tambem na raiz do projeto. 

## Como executar os testes
Seguindos os mesmos passos de como executar a aplicacao, precisamos buildar e rodar a imagem Docker:
```bash
# Dentro da raiz do projeto
docker build . -t stocks-tax-calculator
docker run -it stocks-tax-calculator sh

# com o sh rodando, rodamos os testes
go test ./...

# observaremos a seguinte saida, caso os testes sejam bem suscedidos
?   	github.com/erick-tmr/stocks-tax-calculator/cmd	[no test files]
ok  	github.com/erick-tmr/stocks-tax-calculator/internal/stocks	0.003s
```

## Notas finais
Acabei nao adicionar testes end-to-end testando a aplicacao de ponta a ponta (enviando inputs pro STDIN e lendo a stream do STDOUT) devido a falta de tempo habil e por nao estar familiarizado com este tipo de abordagem em Go, mas acredito ser algo que e possivel de ser implementado, dado algumas horas a mais.
