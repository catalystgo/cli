package main

import (
	"context"

	"github.com/catalystgo/cli/internal/cli"
	log "github.com/catalystgo/logger/cli"
)

func main() {
	if err := cli.Execute(context.Background()); err != nil {
		log.Fatal(err.Error())
	}
}
