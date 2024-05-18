# fullcycle-clean-architecture

# Description of the Challenge (PT/BR)

Olá devs!
Agora é a hora de botar a mão na massa. Para este desafio, você precisará criar o usecase de listagem das orders.
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
  Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.

# How to run

Start the containers using the docker command:

`docker-compose up -d`

Run the database migration:

`make migrate`

To run the project, execute the command:

`cd cmd/ordersystem/ && go run main.go wire_gen.go`

The application will display in the command line that the servers (web, gRPC, and GraphQL) are running and will inform the port, as in the example below:

![img.png](img.png)

## Web Server
To test the web server, execute the files in the `/api` folder of the project which are:

- `create_order.http`
- `list_orders.http`

## GraphQL

To test the GraphQL, access the url http://localhost:8080 and execute the commands to create and list orders:
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

## gRPC

* To test the gRPC, execute the command with the help of `evans` ([see more here](https://evans.syfm.me/)):

```shell
evans -r repl
package pb
service OrderService

```
Then use the command `call CreateOrder` to create an order and `call ListOrders` to list the orders.