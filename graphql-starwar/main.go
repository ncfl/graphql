package main

import (
	"fmt"
	"graphql/graphql-starwar/exec"
	"net/http"

	"github.com/graphql-go/handler"
)

func main() {
	http.Handle("/", handler.New(&handler.Config{
		Schema:     &exec.StarWarsSchema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	}))
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}
