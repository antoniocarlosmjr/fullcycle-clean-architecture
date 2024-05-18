package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	graphQLHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/fullcycle-clean-architecture/config"
	"github.com/fullcycle-clean-architecture/internal/event/handler"
	"github.com/fullcycle-clean-architecture/internal/infra/graph"
	"github.com/fullcycle-clean-architecture/internal/infra/grpc/pb"
	"github.com/fullcycle-clean-architecture/internal/infra/grpc/service"
	"github.com/fullcycle-clean-architecture/internal/infra/web/webserver"
	"github.com/fullcycle-clean-architecture/pkg/events"
	"github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	configs, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	err = eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	if err != nil {
		panic(err)
	}

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrdersUseCase(db)

	// WEB SERVER
	webServer := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)

	webServer.AddHandler("/order", webOrderHandler.Create)
	webServer.AddHandler("/orders", webOrderHandler.List)

	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webServer.Start()

	// GRPC SERVER
	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	// GRAPHQL SERVER
	srv := graphQLHandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	err = http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
	if err != nil {
		panic(err)
	}
}

func getRabbitMQChannel() *amqp091.Channel {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalln("Error to connect rabbitMQ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
