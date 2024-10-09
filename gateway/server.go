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
	"go.mongodb.org/mongo-driver/bson/primitive"
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

		fieldCtx := graphql.GetFieldContext(ctx)
		providedId := fieldCtx.Args["id"].(string)

		if _, err = primitive.ObjectIDFromHex(providedId); err != nil {
			log.Println("Provided incorrect id for token verification.")
			return nil, fmt.Errorf("Please provide a valid id.")
		}
		id, role, err := graph.VerifyToken(token, providedId)

		if err != nil {
			log.Panic(err)
			return nil, err
		}

		ctx = context.WithValue(ctx, "role", role)
		ctx = context.WithValue(ctx, "id", id)
		return next(ctx)

	}

	c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {
		r := ctx.Value("role")
		fmt.Printf("%s %s\n", r, role)

		if r != role.String() {
			log.Println("Unauthorized access.")
			return nil, fmt.Errorf("Only %s is allowed to perform this operation.", role)
		}

		return next(ctx)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "writer", w) // Add the ResponseWriter to context
		srv.ServeHTTP(w, r.WithContext(ctx))
	}))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
