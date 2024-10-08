package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/RushinShah22/e-commerce-micro/gateway/graph"
	"github.com/RushinShah22/e-commerce-micro/gateway/graph/model"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	c := graph.Config{Resolvers: &graph.Resolver{
		UserURL:    "http://users:3000/users",
		ProductURL: "http://products:8000/products",
		OrderURL:   "http://orders:4000/orders",
	}}

	c.Directives.Auth = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		c := graphql.GetOperationContext(ctx)
		header := strings.Split(c.Headers.Get("Authorization"), " ")
		if len(header) != 2 {
			return nil, fmt.Errorf("No authentication token was found.")
		}

		token := header[1]

		err = graph.VerifyToken(token)

		if err != nil {
			return nil, fmt.Errorf("Invalid access token. Please sign in again.")
		}

		return next(ctx)

	}

	c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {

		c := graphql.GetOperationContext(ctx)
		vars := c.Variables

		if vars["role"] != role {
			return nil, fmt.Errorf("Only admins can perform this operation")
		}
		return next(ctx)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
