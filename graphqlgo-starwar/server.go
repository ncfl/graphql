package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

func main() {

	schema := graphql.MustParseSchema(readSchema(), &Resolver{})
	http.Handle("/", http.FileServer(http.Dir("./graphqlgo-starwar/index")))
	http.Handle("/query", &relay.Handler{Schema: schema})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func readSchema() string {
	file, _ := os.Open("graphqlgo-starwar/schema.graphql")
	content, _ := ioutil.ReadAll(file)
	return string(content)
}
