package lxd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (t *Tools) ListImages(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	images, err := t.Client.GetImages()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed listing images: %w", err)
	}

	data, err := json.MarshalIndent(images, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}
