// +build ignore

// This file uses github.com/shurcooL/vfsgen to go generate a fileserver
// implementation which can be used to serve the static files produced by
// the frontend
package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	assets := http.Dir("frontend/public")

	err := vfsgen.Generate(assets, vfsgen.Options{
		Filename:    "frontend/assets.go",
		PackageName: "frontend",
	})
	if err != nil {
		log.Fatal(err)
	}
}
