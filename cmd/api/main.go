package main

import (
	"context"
	"log"
	"wb-challenge/cmd/api/bootstrap"
)

func main() {
	cfg := bootstrap.GetConfigFromEnv()

	logger := log.Default()

	bootstrap.Run(context.Background(), cfg, logger)
}
