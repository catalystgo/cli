package main

import (
	"context"
	"log"

	"github.com/catalystgo/cli/internal/cli"
)

func main() {
	if err := cli.Execute(context.Background()); err != nil {
		log.Fatal(err)
	}
}
