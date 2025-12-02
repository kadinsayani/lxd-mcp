package lxd

import (
	"context"
	"fmt"

	"github.com/canonical/lxd/shared/api"
	"github.com/modelcontextprotocol/go-sdk/mcp"

)

func (t *Tools) StartInstance(ctx context.Context, req *mcp.CallToolRequest, args GetInstanceArgs) (*mcp.CallToolResult, any, error) {
	apiReq := api.InstanceStatePut{
		Action:  "start",
		Timeout: -1,
	}

	op, err := t.Client.UpdateInstanceState(args.Name, apiReq, "")
	if err != nil {
		return nil, nil, fmt.Errorf("Failed starting instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance start: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s started successfully", args.Name)},
		},
	}, nil, nil
}

func (t *Tools) StopInstance(ctx context.Context, req *mcp.CallToolRequest, args StopInstanceArgs) (*mcp.CallToolResult, any, error) {
	apiReq := api.InstanceStatePut{
		Action:  "stop",
		Timeout: -1,
		Force:   args.Force,
	}

	op, err := t.Client.UpdateInstanceState(args.Name, apiReq, "")
	if err != nil {
		return nil, nil, fmt.Errorf("Failed stopping instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance stop: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s stopped successfully", args.Name)},
		},
	}, nil, nil
}

func (t *Tools) RestartInstance(ctx context.Context, req *mcp.CallToolRequest, args RestartInstanceArgs) (*mcp.CallToolResult, any, error) {
	apiReq := api.InstanceStatePut{
		Action:  "restart",
		Timeout: -1,
		Force:   args.Force,
	}

	op, err := t.Client.UpdateInstanceState(args.Name, apiReq, "")
	if err != nil {
		return nil, nil, fmt.Errorf("Failed restarting instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance restart: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s restarted successfully", args.Name)},
		},
	}, nil, nil
}

func (t *Tools) FreezeInstance(ctx context.Context, req *mcp.CallToolRequest, args GetInstanceArgs) (*mcp.CallToolResult, any, error) {
	apiReq := api.InstanceStatePut{
		Action:  "freeze",
		Timeout: -1,
	}

	op, err := t.Client.UpdateInstanceState(args.Name, apiReq, "")
	if err != nil {
		return nil, nil, fmt.Errorf("Failed freezing instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance freeze: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s frozen successfully", args.Name)},
		},
	}, nil, nil
}

func (t *Tools) UnfreezeInstance(ctx context.Context, req *mcp.CallToolRequest, args GetInstanceArgs) (*mcp.CallToolResult, any, error) {
	apiReq := api.InstanceStatePut{
		Action:  "unfreeze",
		Timeout: -1,
	}

	op, err := t.Client.UpdateInstanceState(args.Name, apiReq, "")
	if err != nil {
		return nil, nil, fmt.Errorf("Failed unfreezing instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance unfreeze: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s unfrozen successfully", args.Name)},
		},
	}, nil, nil
}
