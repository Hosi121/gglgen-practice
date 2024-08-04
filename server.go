package main

import (
	"context"
	"example/ent"
	"example/resolver"
	"log"
	"net/http"
	"os"

	"entgo.io/ent/dialect"
	_ "modernc.org/sqlite"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	// Create an ent.Client with in-memory SQLite database.
	entOptions := []ent.Option{}
	entOptions = append(entOptions, ent.Debug())
	client, err := ent.Open(dialect.SQLite, "file::memory:?cache=shared", entOptions...)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// Run the automatic migration tool to create all schema resources.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	ctx := context.Background()
	u, err := client.Todo.
		Create().
		SetName("John").
		SetEmail("hogehoge@gmail.com").
		SetDone(false).
		Save(ctx)
	if err != nil {
		log.Fatalf("failed creating user: %v", err)
	}
	log.Println("user was created: ", u)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(resolver.NewSchema(client))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
