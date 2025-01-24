package main

import (
	"context"
	"github.com/escalopa/grpc-ws/internal/api/chat"

	"github.com/catalystgo/catalystgo"
	"github.com/catalystgo/logger/logger"
)

func main() {
	ctx := context.Background()

	app, err := catalystgo.New()
	if err != nil {
		logger.Fatalf(ctx, "init app: %v", err)
	}

	chatSrv := chat.NewChatService()

	if err := app.Run(chatSrv); err != nil {
		logger.Fatalf(ctx, "app run: %v", err)
	}
}
