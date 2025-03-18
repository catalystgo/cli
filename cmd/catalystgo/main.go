package main

import (
	"context"

	"github.com/catalystgo/cli/internal/cli"
	"github.com/catalystgo/logger/log"
)

func main() {
	if err := cli.Execute(context.Background()); err != nil {
		log.Fatal(err.Error())
	}
}
