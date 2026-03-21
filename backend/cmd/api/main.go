package main

import (
	"context"
	"log"

	"easycpp/backend/internal/app"
)

func main() {
	ctx := context.Background()
	r, container, err := app.NewRouter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}
	defer container.Close()

	addr := container.Config.HTTPAddr
	log.Printf("EasyCpp backend running at %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
