package main

import (
	"authorization_service/internal/app"
	"context"
	"flag"
	"log"
)

func main() {
	flag.Parse()
	ctx := context.Background()

	app, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
