# fullcycle-clean-architecture

* Subir os containers utilizando o comando do docker:
  - `docker-compose up -d`

* Execute a migração do banco de dados:
  - `migrate -path=database/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up`

* Para rodar o projeto, execute o comando:
* `cd cmd/ordersystem/  && go run main.go wire_gen.go`

A aplicação irá exibir em linha de comando que os servidores (web, grpc e graphQL) estão rodando e será informado a porta, como no exemplo a seguir:

![img.png](img.png)

* Para testar o servidor web, execute os arquivos que estão na pasta `/api` do projeto que são:
  * `create_order.http`
  * `list_orders.http`

* Para testar  o GraphQL acesse a url `http://localhost:8080` e execute os comandos para criar e listar os orders:
```graphql
mutation createOrder {
  createOrder(input:{id:"A", Price:30, Tax:2}) {
    id
    Price
    Tax
    FinalPrice
  }
}

 query listOrders {
    orders {
      id
      Tax
      Price
      FinalPrice
   }
 }
```

* Para testar o gRPC, execute o comando com o auxílio do `evans`:

```shell
evans -r repl
```

E faça chamada as seguintes services: `call CreateOrder` ou `call ListOrders`.