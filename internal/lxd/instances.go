package lxd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/canonical/lxd/shared/api"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func (t *Tools) ListInstances(ctx context.Context, req *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	instances, err := t.Client.GetInstances(api.InstanceTypeAny)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed listing instances: %w", err)
	}

	data, err := json.MarshalIndent(instances, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (t *Tools) GetInstance(ctx context.Context, req *mcp.CallToolRequest, args GetInstanceArgs) (*mcp.CallToolResult, any, error) {
	instance, _, err := t.Client.GetInstance(args.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed getting instance: %w", err)
	}

	data, err := json.MarshalIndent(instance, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (t *Tools) GetInstanceState(ctx context.Context, req *mcp.CallToolRequest, args GetInstanceArgs) (*mcp.CallToolResult, any, error) {
	state, _, err := t.Client.GetInstanceState(args.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed getting instance state: %w", err)
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func (t *Tools) CreateInstance(ctx context.Context, req *mcp.CallToolRequest, args CreateInstanceArgs) (*mcp.CallToolResult, any, error) {
	if args.Image == "" {
		return nil, nil, fmt.Errorf("Image parameter is required")
	}

	if args.Type == "" {
		args.Type = "container"
	}

	shouldStart := true
	if args.Start != nil {
		shouldStart = *args.Start
	}

	instanceType := api.InstanceTypeContainer
	if args.Type == "virtual-machine" {
		instanceType = api.InstanceTypeVM
	}

	imgServer, imgAlias, protocol := t.ParseImageSource(args.Image)

	apiReq := api.InstancesPost{
		Name: args.Name,
		Type: instanceType,
		Source: api.InstanceSource{
			Type:     "image",
			Server:   imgServer,
			Protocol: protocol,
			Alias:    imgAlias,
		},
	}

	op, err := t.Client.CreateInstance(apiReq)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed creating instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance creation: %w", err)
	}

	result := fmt.Sprintf("Instance %s created successfully", args.Name)

	if shouldStart {
		startReq := api.InstanceStatePut{
			Action:  "start",
			Timeout: -1,
		}

		op, err := t.Client.UpdateInstanceState(args.Name, startReq, "")
		if err != nil {
			return nil, nil, fmt.Errorf("Instance created but failed starting: %w", err)
		}

		if err := op.Wait(); err != nil {
			return nil, nil, fmt.Errorf("Instance created but failed waiting for start: %w", err)
		}

		result += " and started"
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: result},
		},
	}, nil, nil
}

func (t *Tools) DeleteInstance(ctx context.Context, req *mcp.CallToolRequest, args DeleteInstanceArgs) (*mcp.CallToolResult, any, error) {
	if args.Force {
		state, _, err := t.Client.GetInstanceState(args.Name)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed getting instance state: %w", err)
		}

		if state.Status == "Running" {
			stopReq := api.InstanceStatePut{
				Action:  "stop",
				Timeout: -1,
				Force:   true,
			}

			op, err := t.Client.UpdateInstanceState(args.Name, stopReq, "")
			if err != nil {
				return nil, nil, fmt.Errorf("Failed stopping instance: %w", err)
			}

			if err := op.Wait(); err != nil {
				return nil, nil, fmt.Errorf("Failed waiting for instance stop: %w", err)
			}
		}
	}

	op, err := t.Client.DeleteInstance(args.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed deleting instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance deletion: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s deleted successfully", args.Name)},
		},
	}, nil, nil
}

func (t *Tools) RenameInstance(ctx context.Context, req *mcp.CallToolRequest, args RenameInstanceArgs) (*mcp.CallToolResult, any, error) {
	apiReq := api.InstancePost{
		Name: args.NewName,
	}

	op, err := t.Client.RenameInstance(args.Name, apiReq)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed renaming instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance rename: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s renamed to %s successfully", args.Name, args.NewName)},
		},
	}, nil, nil
}

func (t *Tools) UpdateInstance(ctx context.Context, req *mcp.CallToolRequest, args UpdateInstanceArgs) (*mcp.CallToolResult, any, error) {
	instance, etag, err := t.Client.GetInstance(args.Name)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed getting instance: %w", err)
	}

	if instance.Config == nil {
		instance.Config = make(map[string]string)
	}

	for key, value := range args.Config {
		instance.Config[key] = value
	}

	op, err := t.Client.UpdateInstance(args.Name, instance.Writable(), etag)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed updating instance: %w", err)
	}

	if err := op.Wait(); err != nil {
		return nil, nil, fmt.Errorf("Failed waiting for instance update: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Instance %s updated successfully", args.Name)},
		},
	}, nil, nil
}
