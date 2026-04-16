package main

import (
	"context"
	"flag"
	"github.com/Coldwws/auth_service/internal/app"
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
