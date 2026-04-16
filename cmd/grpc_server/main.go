package main

import (
	"authorization_service/internal/app"
	"context"
	"flag"
	"log"
)

var logLevel string

func init() {
	flag.StringVar(&logLevel, "l", "info", "log level")
}
func main() {
	flag.Parse()
	ctx := context.Background()

	app, err := app.NewApp(ctx, logLevel)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
