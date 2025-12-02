package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kadinsayani/lxd-mcp/internal/lxd"
	"github.com/kadinsayani/lxd-mcp/internal/server"
)

func main() {
	tools, err := lxd.NewTools()
	if err != nil {
		log.Fatalf("Failed creating LXD tools: %v", err)
	}

	srv := server.Setup(tools)

	if err := srv.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
