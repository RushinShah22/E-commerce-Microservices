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
)

const defaultPort = "8080"

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	uURL := "http://users:3000/users"
	pURL := "http://products:8000/products"
	oURL := "http://orders:4000/orders"
	c := graph.Config{Resolvers: &graph.Resolver{
		UserURL:    uURL,
		ProductURL: pURL,
		OrderURL:   oURL,
	}}

	c.Directives.Auth = func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {

		ct := graphql.GetOperationContext(ctx)
		header := strings.Split(ct.Headers.Get("Authorization"), " ")
		if len(header) != 2 {
			err = fmt.Errorf("No bearer token was found")
			return
		}

		token := header[1]

		id, role, err := graph.VerifyToken(token)

		if err != nil {
			log.Println(err)
			return
		}

		resp, err := http.Get(uURL + "/" + id)

		if err != nil {
			log.Println(err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			log.Println("Tried to access a admin path using fake user id.")
			err = fmt.Errorf("Your session has expired. Please sign in again.")
			return
		}

		ctx = context.WithValue(ctx, "role", role)
		ctx = context.WithValue(ctx, "id", id)

		return next(ctx)
	}

	c.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (res interface{}, err error) {
		r := ctx.Value("role")

		if r != role.String() {
			log.Println("Unauthorized access.")
			return struct{}{}, fmt.Errorf("Only %s is allowed to perform this operation.", role)
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
