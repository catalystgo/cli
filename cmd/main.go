package main

import (
	"context"
	"log"
	"time"

	"github.com/catalystgo/cli/internal/cli"
)

const timeoutDuration = 10 * time.Second

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	if err := cli.Execute(ctx); err != nil {
		log.Fatal(err)
	}
}
