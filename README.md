## SWAPI Challenge (OUTDATED - Não use esse projeto como referencia para nada)

#### Requisitos:
* A API deve ser REST
* Para cada planeta, os seguintes dados devem ser obtidos do banco de dados da aplicação, sendo inserido manualmente:
  * Nome
  * Clima
  * Terreno
  * Quantidade de aparições em filmes (que podem ser obtidas pela API pública do Star Wars:  https://swapi.co/)

#### Funcionalidades: 
* Adicionar um planeta (com nome, clima e terreno)
* Listar planetas
* Buscar por nome
* Buscar por ID
* Remover planeta

#### Environment vars
* MONGO_URL - URL de acesso ao mongodb (ex: mongodb://<user>:<pass>@<host>:<port>)
* SERVER_ADDR - Onde a App vai fazer bind (ex ":8080")
* API_RATE_LIMIT - Quantidade de requisicoes por segundo que a App ira atender (ex: 20) (default: 100) 
* CB_REQ_COUNT - Quantidade de requests para considerar fechar o circuito (ex: 40) (default: 3)
* CB_REQ_FAIL_RATIO - razao (total falhas):(total de request) (ex: 1.5) (default: 0.5)

### FAQ

##### Como eu vejo a App funcionando?
No diretorio raiz:
`docker-compose up`

`Os exemplos abaixo sao chamadas usando o httpie, mas sao facilmente tarduziadas para o uso do bom e velho curl`
###### Adicionar um planeta
`* http POST localhost:8000 climate=seco name=Tatooine terrain=terreno`

Retornos:

    * 200 : Sucesso
    * 400 : Erros de validacao de input 
    * 500 : Erro inesperado
    
###### Listar planetas
`* http GET localhost:8000/planets`
Retornos:

    * 200 : Sucesso
    * 500 : Erro inesperado
###### Buscar por nome
`* http GET localhost:8000/planets?name=<NOME_DO_PLANETA_AQUI>`

Retornos:

    * 200 : Sucesso
    * 400 : Erros de validacao de input 
    * 500 : Erro inesperado
###### Buscar por ID
`* http GET localhost:8000/planets?id=<IDENTIFICADOR_DO_PLANETA_AQUI>`

Retornos:

    * 200 : Sucesso
    * 400 : Erros de validacao de input 
    * 500 : Erro inesperado
###### Remover planeta
`* http DELETE localhost:8000/planets?id=<IDENTIFICADOR_DO_PLANETA_AQUI>`

Retornos:

    * 200 : Sucesso
    * 400 : Erros de validacao de input 
    * 500 : Erro inesperado
##### Como eu rodo a App standalone sem usar docker-compose?
No diretorio raiz:
 `MONGO_URL=mongodb://<user>:<pass>@<host>:<port> SERVER_ADDR=:8080 go run main.go`

#### Sobre o gerenciamento de dependencias, por que Go Dep?
Foi escolhido a ferramenta [dep](https://golang.github.io/dep/) para o gerenciamento de dependencias, pois sera a ferramenta incorporada na linguagem num futuro proximo.
Isso vai evitar reescritas pra adapatacao.

#### Por que voce versionou o diretorio `/vendor`?
O repositorio fica maior e os diffs sao mais complicados de lidar, porem para esse desafio em especifico eu quis garantir 2 coisas:
* O usuario desse projeto nao vai precisar executar o `dep ensure` para buildar o programa.
* Garantir que o codigo vai reproduzir o mesmo comportamento no build por estar desacoplado de repositorios externos (quase como um cache)

#### Durante o build vc manteve o CGO, pq?
A imagem base para o build e a de execucao sao a mesma distro, nao precisei me preocupar com falhas na linkagem dinamica.

#### Para que foi utilizado um rate limiter
Para controle de vazao

#### Para que foi utilizado um circuit breaker
Para nao ter problemas de performance da nossa App originados pelo meio externo

#### O que poderia ter a mais que nao teve?
###### Para ser de facil manutencao:
* Commits separados com a progressao do desenvolvimento, facilitaria um tracking no futuro
* Swagger automatizado
* Externalizar o base path como variavel de ambiente
###### Para qualidade:
* Testes unitarios (Essa parte eh facil, mas demora pois tem muito boilerplate, entao preferi focar em features mais interessante para o desafio)
* Teste de aceitacao (Acho interessante tbm a abordagem do BDD, por ex: https://github.com/DATA-DOG/godog/)
